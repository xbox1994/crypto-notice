vim crypto-notice.service
systemctl daemon-reload
systemctl enable crypto-notice.service
systemctl start crypto-notice.service
systemctl status crypto-notice.service