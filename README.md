# alsamixer2mqtt

This project allows you to bridge ALSA sound controls (e.g., volume levels) with MQTT, enabling integration with Home Assistant or other systems. With this application, you can monitor sound levels and adjust them remotely using MQTT commands.

## Installation

Run as any user on your RaspberryPi:

```bash
bash <(wget -qO- https://raw.githubusercontent.com/sokoli-media/alsamixer2mqtt/main/install.sh)
```

## HASS configuration

Add to `configuration.yaml`:

```yaml
mqtt:
  - number:
      name: Friendly name
      object_id: not_friendly_name_used_for_entity_id
      icon: mdi:account-voice
      command_topic: topic_you_set_in_the_installation_process
      state_topic: topic_you_set_in_the_installation_process
      retain: true
      unit_of_measurement: '%'
      min: 0
      max: 100
      step: 1
```
