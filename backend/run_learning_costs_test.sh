#!/bin/bash

# Wechseln ins Verzeichnis des Tests
cd /data/dev/bamort/backend

# Test ausführen mit detaillierter Ausgabe
go test -v ./gsmaster -run TestCalculateDetailedSkillLearningCostForHexer
