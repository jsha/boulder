#!/usr/bin/env python2.7
"""Receive a single SMTP message and print it."""
import smtpd
import asyncore
import sys

class SingleDebuggingServer(smtpd.SMTPServer):
    def __init__(self, local):
        smtpd.SMTPServer.__init__(self, local, None)

    def process_message(self, peer, mailfrom, rcpttos, data):
        print data
        sys.exit(0)

server = SingleDebuggingServer(('127.0.0.1', 1025))

asyncore.loop(timeout=5, count=1)
