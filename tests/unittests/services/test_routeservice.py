import unittest
import unittest.mock as mock
import inspect
from transiter import models
from transiter.services import routeservice


class TestRouteService(unittest.TestCase):

    __module_being_tested__ = routeservice

    SYSTEM_ID = '1'

    ROUTE_ONE_PK = 2
    ROUTE_ONE_ID = '3'
    ROUTE_ONE_STATUS = routeservice.Status.PLANNED_SERVICE_CHANGE

    ROUTE_TWO_PK = 4
    ROUTE_TWO_ID = '5'
    ROUTE_TWO_STATUS = routeservice.Status.GOOD_SERVICE

    RAW_FREQUENCY = 700

    @classmethod
    def _getImportedModulesDict(cls, _caches={}):
        if len(_caches) == 0:
            _caches['module_cache'] = {}
            for name, obj in cls.__module_being_tested__.__dict__.items():
                if name[:2] == '__':
                    continue
                if inspect.ismodule(obj):
                    _caches['module_cache'][obj.__name__] = name
        return _caches['module_cache']

    def mockModuleVariable(self, variable_name):
        patcher = mock.patch.object(
            self.__module_being_tested__,
            variable_name
        )
        mocked_module = patcher.start()
        self.addCleanup(patcher.stop)
        return mocked_module

    def mockImportedModule(self, imported_module):
        module_global_name = getattr(imported_module, '__name__')
        module_local_name = self._getImportedModulesDict()[module_global_name]
        return self.mockModuleVariable(module_local_name)

    def setUp(self):
        self.route_one = models.Route()
        self.route_one.id = self.ROUTE_ONE_ID
        self.route_one.pk = self.ROUTE_ONE_PK
        self.route_one.service_patterns = []
        self.route_one.alerts = []

        self.route_two = models.Route()
        self.route_two.id = self.ROUTE_TWO_ID
        self.route_two.pk = self.ROUTE_TWO_PK

        self.routedam = self.mockImportedModule(routeservice.routedam)
        self.systemdam = self.mockImportedModule(routeservice.systemdam)

    @mock.patch.object(routeservice, '_construct_route_pk_to_status_map')
    def test_list_all_in_system(self, _construct_route_pk_to_status_map):
        """[Route service] List all routes in a system"""

        _construct_route_pk_to_status_map.return_value = {
            self.ROUTE_ONE_PK: self.ROUTE_ONE_STATUS,
            self.ROUTE_TWO_PK: self.ROUTE_TWO_STATUS
        }
        self.routedam.list_all_in_system.return_value = [
            self.route_one,
            self.route_two
        ]
        self.systemdam.get_by_id.return_value = models.System()

        expected = [
            {
                **self.route_one.short_repr(),
                'status': self.ROUTE_ONE_STATUS
            },
            {
                **self.route_two.short_repr(),
                'status': self.ROUTE_TWO_STATUS
            }
        ]

        actual = routeservice.list_all_in_system(self.SYSTEM_ID)

        self.assertEqual(actual, expected)

        self.routedam.list_all_in_system.assert_called_once_with(self.SYSTEM_ID)

    @mock.patch.object(routeservice, '_construct_route_status')
    def test_get_in_system_by_id(self, _construct_route_status):
        """[Route service] Get a specific route in a system"""

        _construct_route_status.return_value = self.ROUTE_ONE_STATUS
        self.routedam.get_in_system_by_id.return_value = self.route_one
        self.routedam.calculate_frequency.return_value = self.RAW_FREQUENCY

        expected = {
            **self.route_one.short_repr(),
            'frequency': int(self.RAW_FREQUENCY/6)/10,
            'status': self.ROUTE_ONE_STATUS,
            'alerts': [],
            'service_maps': []
        }

        actual = routeservice.get_in_system_by_id(
            self.SYSTEM_ID,
            self.ROUTE_ONE_ID
        )

        self.assertDictEqual(actual, expected)

        self.routedam.get_in_system_by_id.assert_called_once_with(
            self.SYSTEM_ID,
            self.ROUTE_ONE_ID
        )

    @mock.patch.object(routeservice, '_construct_route_pk_to_status_map')
    def test_construct_route_status(self, _construct_route_pk_to_status_map):

        _construct_route_pk_to_status_map.return_value = {
            self.ROUTE_ONE_PK: self.ROUTE_ONE_STATUS
        }

        self.assertEqual(
            self.ROUTE_ONE_STATUS,
            routeservice._construct_route_status(self.ROUTE_ONE_PK)
        )




