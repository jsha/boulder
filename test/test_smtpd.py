#!/usr/bin/env python2.7
"""Receive a single SMTP message and print it."""
import smtpd
import asyncore

server = smtpd.DebuggingServer(('127.0.0.1', 1025), None)

asyncore.loop()
