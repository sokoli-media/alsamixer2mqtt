# alsamixer2mqtt

This project allows you to bridge ALSA sound controls (e.g., volume levels) with MQTT, enabling integration with Home Assistant or other systems. With this application, you can monitor sound levels and adjust them remotely using MQTT commands.

## Installation

```bash
bash <(wget -qO- https://raw.githubusercontent.com/sokoli-media/alsamixer2mqtt/main/install.sh)
```

## HASS configuration

Add to `configuration.yaml`:

```yaml
mqtt:
  - number:
      name: any_name_you_may_think_about
      command_topic: topic_you_set_in_the_installation_process
      state_topic: topic_you_set_in_the_installation_process
      step: 1
      min: 0
      max: 100
```
