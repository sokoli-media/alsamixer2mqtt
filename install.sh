#!/bin/bash

# Function to prompt for environment variables
prompt_env_variable() {
    local var_name=$1
    local default_value=$2
    local input

    read -p "Enter value for $var_name (default: $default_value): " input
    echo "${input:-$default_value}"
}

# Step 1: Ask for installation directory
INSTALL_DIR="/opt/alsamixer2mqtt"
read -p "Enter installation directory (default: $INSTALL_DIR): " input
INSTALL_DIR="${input:-$INSTALL_DIR}"

# Step 2: Check if .env file already exists
ENV_FILE="$INSTALL_DIR/.env"
if [ -f "$ENV_FILE" ]; then
    echo "Environment file already exists at $ENV_FILE."
    echo "Skipping environment variable setup."
else
    # Step 3: Prompt for environment variables if .env file does not exist
    echo "No existing .env file found. Prompting for environment variables..."
    MQTT_BROKER=$(prompt_env_variable "MQTT_BROKER" "tcp://localhost:1883")
    MQTT_CLIENT_ID=$(prompt_env_variable "MQTT_CLIENT_ID" "alsamixer2mqtt")
    MQTT_USERNAME=$(prompt_env_variable "MQTT_USERNAME" "")
    MQTT_PASSWORD=$(prompt_env_variable "MQTT_PASSWORD" "")
    ALSA_DEVICE=$(prompt_env_variable "ALSA_DEVICE" "default")
    ALSA_CONTROL=$(prompt_env_variable "ALSA_CONTROL" "Master")
    STATE_TOPIC=$(prompt_env_variable "STATE_TOPIC" "")
    SET_TOPIC=$(prompt_env_variable "SET_TOPIC" "")

    # Step 4: Create the .env file with the environment variables
    echo "Creating .env file..."
    sudo mkdir -p "$INSTALL_DIR"
    sudo tee "$ENV_FILE" > /dev/null <<EOF
MQTT_BROKER=$MQTT_BROKER
MQTT_CLIENT_ID=$MQTT_CLIENT_ID
MQTT_USERNAME=$MQTT_USERNAME
MQTT_PASSWORD=$MQTT_PASSWORD
ALSA_DEVICE=$ALSA_DEVICE
ALSA_CONTROL=$ALSA_CONTROL
STATE_TOPIC=$STATE_TOPIC
SET_TOPIC=$SET_TOPIC
EOF
fi

# Step 5: Download the latest artifact
echo "Downloading the latest artifact from the GitHub repository..."
LATEST_ARTIFACT_URL="https://github.com/sokoli-media/alsamixer2mqtt/releases/latest/download/alsamixer2mqtt"
wget -O alsamixer2mqtt "$LATEST_ARTIFACT_URL" || {
    echo "Error: Failed to download artifact."
    exit 1
}

# Step 6: Install the application
echo "Installing application to $INSTALL_DIR..."
sudo mkdir -p "$INSTALL_DIR"
sudo mv alsamixer2mqtt "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/alsamixer2mqtt"

# Step 7: Create the systemd service file
SERVICE_FILE="/etc/systemd/system/alsamixer2mqtt.service"
echo "Creating systemd service file at $SERVICE_FILE..."
sudo tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=ALSAMixer to MQTT Service
After=network.target

[Service]
Type=simple
EnvironmentFile=$ENV_FILE
ExecStart=$INSTALL_DIR/alsamixer2mqtt
Restart=on-failure
User=$(whoami)

[Install]
WantedBy=multi-user.target
EOF

if systemctl is-active --quiet alsamixer2mqtt.service; then
  echo "Restarting the service..."
  sudo systemctl daemon-reload
  sudo systemctl restart alsamixer2mqtt.service
else
  # Step 8: Reload systemd, enable and start the service
  echo "Enabling and starting the service..."
  sudo systemctl daemon-reload
  sudo systemctl enable alsamixer2mqtt.service
  sudo systemctl start alsamixer2mqtt.service
fi

echo "Installation complete. The service is now running."

# Display instructions for accessing logs
echo ""
echo "You can now check the service status and logs using the following commands:"
echo ""
echo "To check the status of the service:"
echo "  sudo systemctl status alsamixer2mqtt"
echo ""
echo "To view logs from the service using journalctl:"
echo "  sudo journalctl -u alsamixer2mqtt -f"
echo ""
echo "Logs will be displayed in real-time with the -f flag."
echo ""
echo "If you need to stop the service, use:"
echo "  sudo systemctl stop alsamixer2mqtt"
echo ""
echo "To restart the service:"
echo "  sudo systemctl restart alsamixer2mqtt"
