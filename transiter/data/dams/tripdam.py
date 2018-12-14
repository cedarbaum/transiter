from transiter import models
from transiter.data import database


def list_all_in_route_by_pk(route_pk):
    session = database.get_session()
    query = (
        session.query(models.Trip)
        .filter(models.Trip.route_pk == route_pk)
    )
    for row in query:
        yield row


def list_all_in_routes(system_id, route_ids):
    session = database.get_session()
    query = session.query(models.Trip) \
        .join(models.Route, models.Route.pk == models.Trip.route_pk) \
        .filter(models.Route.system_id == system_id) \
        .filter(models.Route.id.in_(route_ids))
    for row in query:
        yield row


# TODO: make this call the one above
# TODO: are both of these actually used?
def list_all_in_route(system_id, route_id):
    session = database.get_session()
    query = session.query(models.Trip) \
        .join(models.Route, models.Route.pk == models.Trip.route_pk) \
        .filter(models.Route.system_id == system_id) \
        .filter(models.Route.id == route_id)
    for row in query:
        yield row


def get_in_route_by_id(system_id, route_id, trip_id):
    session = database.get_session()
    return session.query(models.Trip) \
        .join(models.Route, models.Route.pk == models.Trip.route_pk) \
        .filter(models.Route.system_id == system_id) \
        .filter(models.Route.id == route_id) \
        .filter(models.Trip.id == trip_id) \
        .one_or_none()


def get_trip_pk_to_future_stop_events_map(trip_pks):
    session = database.get_session()
    query = (
        session.query(models.StopTimeUpdate)
        .filter(models.StopTimeUpdate.trip_pk.in_(trip_pks))
        .filter(models.StopTimeUpdate.future == True)
        .order_by(models.StopTimeUpdate.stop_sequence)
        .all()
    )
    result = {trip_pk: [] for trip_pk in trip_pks}
    for row in query:
        result[row.trip_pk].append(row)
    return result
