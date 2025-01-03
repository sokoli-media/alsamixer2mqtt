# alsamixer2mqtt

Fun fact: this whole project was generated using ChatGPT to learn
how far can we succeed with it. If there are any random commits
after the first one, it means that ChatGPT failed (or I asked
it to generate more changes).

Content below has also been generated automatically.

---

This project allows you to bridge ALSA sound controls (e.g., volume levels) with MQTT, enabling integration with Home Assistant or other systems. With this application, you can monitor sound levels and adjust them remotely using MQTT commands.

---

## Features

- Publishes sound levels (e.g., volume) as MQTT topics.
- Allows control of ALSA sound levels (e.g., "Master") via MQTT.
- Easily configurable via environment variables.
- Designed for Raspberry Pi or other Linux-based systems.

---

## Installation

### Prerequisites
1. **Linux System**: This application is designed for Linux systems with ALSA support.
2. **Dependencies**: Ensure `wget`, `systemd`, and ALSA are installed.
3. **MQTT Broker**: You will need access to an MQTT broker (e.g., Mosquitto).

---

### One-Liner Installation Command

Run the following command to install and set up the service:

```bash
bash <(wget -qO- https://raw.githubusercontent.com/sokoli-media/alsamixer2mqtt/main/install.sh)