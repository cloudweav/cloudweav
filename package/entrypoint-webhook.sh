#!/bin/bash
set -e

exec tini -- cloudweav-webhook "${@}"
