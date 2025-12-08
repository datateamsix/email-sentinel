#!/bin/sh
set -e

# Post-installation script for email-sentinel

echo "Email Sentinel installed successfully!"
echo ""
echo "To get started:"
echo "  1. Run: email-sentinel init"
echo "  2. Add a filter: email-sentinel filter add"
echo "  3. Start monitoring: email-sentinel start --tray"
echo ""
echo "To enable auto-start on login:"
echo "  systemctl --user enable email-sentinel"
echo "  systemctl --user start email-sentinel"
echo ""
echo "Documentation: https://github.com/datateamsix/email-sentinel"
