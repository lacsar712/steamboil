# Steamboil

Power plant boiler drum and combustion coordination service with embedded operator HMI.

## Features

- Boiler steam pressure and feedwater regulation
- Combustion purge, ignition, and load ramp coordination
- Drum level, swell/shrink detection, and carryover monitoring
- FSM-driven plant states with hook chains
- Process clock timing windows (purge, ignition delay, drum swell settle)
- Interlock leases with defer release and safety gates
- Snapshot store with revision tracking and journal persistence
- Embedded web HMI (~30% frontend)

## Run locally

```bash
go build -o steamboil ./cmd/steamboil
./steamboil
```

Open http://localhost:8080 for the HMI. API lives under `/api/`.

## Docker

```bash
sh build_benzhi_docker.sh
docker run --rm -p 8080:8080 steamboil:benzhi
```

## Module

`github.com/lacsar712/steamboil`
