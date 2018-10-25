from ..database import connection
from ..database import models
from transiter.database.daos import direction_name_dao, feed_dao, route_dao, stop_dao, station_dao, system_dao
from . import exceptions
from transiter.utils import gtfsstaticutil
from ..utils import routelistutil
import csv
import os
from ..utils import jsonutil
from transiter.utils import linksutil
from transiter.utils import servicepatternmanager

@connection.unit_of_work
def list_all():
    """
    List all installed systems.
    :return: A list of short representation of systems
    """
    response = []
    for system in system_dao.list_all():
        system_response = system.short_repr()
        system_response.update({
            'href': linksutil.SystemEntityLink(system)
        })
        response.append(system_response)
    return response
    [system.short_repr() for system in system_dao.list_all()]


@connection.unit_of_work
def get_by_id(system_id):
    system = system_dao.get_by_id(system_id)
    if system is None:
        raise exceptions.IdNotFoundError
    response = system.short_repr()
    response.update({
        "stops": {
            "count": system_dao.count_stops_in_system(system_id),
            "href": linksutil.StopsInSystemIndexLink(system)
        },
        "stations": {
            "count": system_dao.count_stations_in_system(system_id),
            "href": "NI"
        },
        "routes": {
            "count": system_dao.count_routes_in_system(system_id),
            "href": linksutil.RoutesInSystemIndexLink(system)
        },
        "feeds": {
            "count": system_dao.count_feeds_in_system(system_id),
            "href": linksutil.FeedsInSystemIndexLink(system)
        }
    })
    return response


@connection.unit_of_work
def install(system_id):
    if system_dao.get_by_id(system_id) is not None:
        return False

    system = system_dao.create()
    system.system_id = system_id

    _import_static_data(system)
    return True


@connection.unit_of_work
def delete_by_id(system_id):
    deleted = system_dao.delete_by_id(system_id)
    if not deleted:
        raise exceptions.IdNotFoundError
    return True



def _read_csv_file(file_path):
    with open(file_path, mode='r') as csv_file:
        csv_reader = csv.DictReader(csv_file)
        yield from csv_reader


def _import_static_data(system):
    system_base_dir = os.path.join(
        os.path.dirname(__file__),
        '../systems',
        system.system_id
        )
    agency_data_dir = os.path.join(system_base_dir, 'agencydata')
    custom_data_dir = os.path.join(system_base_dir, 'customdata')
    print(agency_data_dir)

    gtfs_static_parser = gtfsstaticutil.GtfsStaticParser()
    gtfs_static_parser.parse_from_directory(agency_data_dir)

    trip_id_to_trip = gtfs_static_parser.trip_id_to_trip
    for trip in trip_id_to_trip.values():
        print(trip.stop_ids)
        break
    exit()


    routes_data_file = os.path.join(agency_data_dir, 'routes.txt')
    routes_by_route_id = {}
    for row in _read_csv_file(routes_data_file):
        route = route_dao.create()
        route.route_id=row['route_id']
        route.color = row['route_color']
        route.timetable_url = row['route_url']
        route.short_name = row['route_short_name']
        route.long_name = row['route_long_name']
        route.description = row['route_desc']
        route.system = system
        routes_by_route_id[route.route_id] = route

    stops_data_file = os.path.join(agency_data_dir, 'stops.txt')
    station_sets_by_stop_id = {}
    stops_by_stop_id = {}
    for row in _read_csv_file(stops_data_file):
        # Note, this may be NYC subway specific logic: if so extract the
        # it to the NYC subway module
        stop_id = row['stop_id']
        if stop_id[-1] == 'N' or stop_id[-1] == 'S':
            continue

        stop = stop_dao.create()
        stop.stop_id=row['stop_id']
        stop.name = row['stop_name']
        stop.longitude = row['stop_lon']
        stop.lattitude = row['stop_lat']
        stop.system = system
        station_sets_by_stop_id[stop_id] = {stop_id}
        stops_by_stop_id[stop_id] = stop

    transfers_data_file = os.path.join(agency_data_dir, 'transfers.txt')
    for row in _read_csv_file(transfers_data_file):
        stop_id_1 = row['from_stop_id']
        stop_id_2 = row['to_stop_id']
        if stop_id_1 == stop_id_2:
            continue

        updated_station_set = station_sets_by_stop_id[stop_id_1].union(
            station_sets_by_stop_id[stop_id_2])
        for stop_id in updated_station_set:
            station_sets_by_stop_id[stop_id] = updated_station_set

    for station_set in station_sets_by_stop_id.values():
        if len(station_set) == 0:
            continue

        station = station_dao.create()
        for stop_id in station_set:
            stops_by_stop_id[stop_id].station = station
        station.system = system
        station_set.clear()

    stop_times_data_file = os.path.join(agency_data_dir, 'stop_times.txt')
    servicepatternmanager.generate_service_patterns_from_gtfs_static_data(
        system,
        stop_times_data_file,
        routes_by_route_id,
        stops_by_stop_id)
    """
    route_lists = routelistutil.construct_route_lists_from_stop_times_file(
        system,
        stop_times_data_file
    )
    for (route_id, route_list) in route_lists.items():
        route = routes_by_route_id[route_id]
        position = 0
        for stop_id in route_list:
            # TODO put this in a DAO when ServicePattern is a thing
            route_list_entry = models.RouteListEntry()
            session = connection.get_session()
            session.add(route_list_entry)
            route_list_entry.route = route
            route_list_entry.stop = stops_by_stop_id[stop_id]
            route_list_entry.position = position
            position += 1
    """

    # The following two data imports are definitely custom logic, though
    # custom to the program rather than the NYC subway
    # Default other option: N=north, S=south ;)
    # Also direction name exceptions is a bad name
    direction_names_data_file = os.path.join(
        custom_data_dir,
        'direction_names.csv'
        )
    for row in _read_csv_file(direction_names_data_file):
        stop_id = row['stop_id']
        north = direction_name_dao.create()
        north.name = row['north_direction_name']
        north.track = None
        north.direction = 'N'
        north.stop = stops_by_stop_id[stop_id]

        south = direction_name_dao.create()
        south.name = row['south_direction_name']
        south.track = None
        south.direction = 'S'
        south.stop = stops_by_stop_id[stop_id]











    direction_name_exceptions_data_file = os.path.join(
        custom_data_dir,
        'direction_name_exceptions.csv'
        )
    for row in _read_csv_file(direction_name_exceptions_data_file):
        direction = direction_name_dao.create()

        direction.name = row['name']
        direction.track = row['track']
        direction.direction = row['direction']
        direction.stop = stops_by_stop_id[row['stop_id']]

    feeds_data_file = os.path.join(custom_data_dir, 'feeds.csv')
    for row in _read_csv_file(feeds_data_file):
        feed = feed_dao.create()
        feed.system = system
        feed.feed_id = row['feed_id']
        feed.url = row['url']
        feed.parser_module = row['parser_module']
        feed.parser_function = row['parser_function']
