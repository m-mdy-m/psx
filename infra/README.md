
Multiple Dockerfile variants available:

### Standard: `Dockerfile`
Multi-stage build with Alpine base.

```bash
docker buildx build  -t psx:latest .
sudo docker run -it --rm -v "$(pwd)":/project -w /project psx:latest check
```

**Size:** ~20MB

### Alpine: `infra/Dockerfile.alpine`
Optimized Alpine-based image.

```bash
docker buildx build -t psx:alpine -f infra/Dockerfile.alpine .
docker run -it --rm -v $(pwd):/project -w /project psx:alpine check
```

**Size:** ~17MB

### Scratch: `infra/Dockerfile.scratch`
Minimal scratch-based image (smallest).

```bash
docker buildx build -t psx:scratch -f infra/Dockerfile.scratch .
docker run -it --rm -v $(pwd):/project -w /project psx:scratch check
```

**Size:** ~5MB

### Docker Compose

```bash
# Build all variants
docker-compose build

# Run standard version
docker-compose up psx

# Run Alpine version
docker-compose up psx-alpine

# Run Scratch version
docker-compose up psx-scratch
```