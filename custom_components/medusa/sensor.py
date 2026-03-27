"""Platform for sensor integration."""
from homeassistant.components.sensor import SensorEntity, SensorDeviceClass, SensorStateClass
from homeassistant.const import (
    UnitOfTemperature,
    PERCENTAGE,
    LIGHT_LUX,
    UnitOfElectricPotential,
    UnitOfPressure,
    UnitOfLength
)
from homeassistant.core import callback
from homeassistant.helpers.dispatcher import async_dispatcher_connect

from .const import (
    DOMAIN, 
    SIGNAL_MEDUSA_UPDATE, 
    ACTION_TEMP, 
    ACTION_LIGHT, 
    ACTION_VOLT, 
    ACTION_GAS, 
    ACTION_PRESSURE, 
    ACTION_ALTITUDE
)

async def async_setup_platform(hass, config, async_add_entities, discovery_info=None):
    """Set up the sensor platform."""
    if discovery_info is None:
        return

    hub = hass.data[DOMAIN]["hub"]
    config_data = hass.data[DOMAIN]["config"]
    
    entities = []
    
    for board in config_data.get("Boards", []):
        room = board.get("Room", "unknown")
        name = board.get("Name", "unknown")
        actions = board.get("Actions", [])
        
        if ACTION_TEMP in actions:
            entities.append(MedusaSensor(hub, room, name, "temp", SensorDeviceClass.TEMPERATURE, UnitOfTemperature.FAHRENHEIT))
            entities.append(MedusaSensor(hub, room, name, "humidity", SensorDeviceClass.HUMIDITY, PERCENTAGE))
            
        if ACTION_LIGHT in actions:
            entities.append(MedusaSensor(hub, room, name, "light", SensorDeviceClass.ILLUMINANCE, LIGHT_LUX))
            
        if ACTION_VOLT in actions:
            entities.append(MedusaSensor(hub, room, name, "volt", SensorDeviceClass.VOLTAGE, UnitOfElectricPotential.VOLT))
            
        if ACTION_GAS in actions:
            entities.append(MedusaSensor(hub, room, name, "gas", None, "kOhm"))
            
        if ACTION_PRESSURE in actions:
            entities.append(MedusaSensor(hub, room, name, "pressure", SensorDeviceClass.ATMOSPHERIC_PRESSURE, UnitOfPressure.PA))
            
        if ACTION_ALTITUDE in actions:
            entities.append(MedusaSensor(hub, room, name, "altitude", SensorDeviceClass.DISTANCE, UnitOfLength.METERS))

    async_add_entities(entities)


class MedusaSensor(SensorEntity):
    """Representation of a Medusa Sensor."""

    def __init__(self, hub, room, board_name, sensor_type, device_class, unit):
        """Initialize the sensor."""
        self._hub = hub
        self._room = room
        self._board_name = board_name
        self._sensor_type = sensor_type
        
        self._attr_name = f"{room} {board_name} {sensor_type}"
        self._attr_unique_id = f"{room}_{board_name}_{sensor_type}"
        self._attr_device_class = device_class
        self._attr_native_unit_of_measurement = unit
        self._attr_state_class = SensorStateClass.MEASUREMENT
        self._state = None
        
        self._attr_device_info = {
            "identifiers": {(DOMAIN, f"{room}_{board_name}")},
            "name": f"{room}_{board_name}",
            "manufacturer": "Medusa",
            "suggested_area": room,
        }

    @property
    def native_value(self):
        """Return the state of the sensor."""
        return self._state

    async def async_added_to_hass(self):
        """Register callbacks."""
        self.async_on_remove(
            async_dispatcher_connect(
                self.hass,
                SIGNAL_MEDUSA_UPDATE.format(self._room, self._board_name),
                self._handle_update,
            )
        )

    @callback
    def _handle_update(self, data):
        """Handle updated data from the Medusa Hub."""
        if self._sensor_type in data:
            self._state = data[self._sensor_type]
            self.async_write_ha_state()
