#!/bin/sh
set -e

# Pre-removal script for email-sentinel

# Stop and disable service if running
if systemctl --user is-active --quiet email-sentinel 2>/dev/null; then
    echo "Stopping email-sentinel service..."
    systemctl --user stop email-sentinel || true
fi

if systemctl --user is-enabled --quiet email-sentinel 2>/dev/null; then
    echo "Disabling email-sentinel service..."
    systemctl --user disable email-sentinel || true
fi
