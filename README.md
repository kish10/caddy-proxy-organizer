# caddy-proxy-organizer
Organize Caddy reverse-proxy service with Docker

## Approaches

### 1. Docker service/container discovery through the Docker socket on the Docker Bridge network

Pros:
- Simple

Cons:
- Limited set up.
    - All the services run on a single machine

### 2. Docker service/container discovery through the Docker overlay network

Pros:
- Compared to approach 1, can ran services on multiple machines

Cons:
- Compared to approach 4, need to manage service set up & service duplication (for load balancing/scaling) manually.

### 3. Docker service/container discovery through Consul (& Consul template)

Pros:
- Compared to approach 1 & 2, take you closer to less reliance on Docker for service discovery & management.
  - Although still using either a Docker bridge network or Docker overlay network, it is a step closer to using more robust service orchestration strategies such as Nomad or Kubernates, see step 4.

Cons:
- Compared to approach 1, have perform the extra step of setting up the Consul server.

### 4. Using Nomad + Consul (& Consul template)

(The reccommended approach)

Pros:
- Simple well written library that takes care of alot tedious work of set up & configuration, and positions you to easily leverage more sophisticated service management strategies. 

Cons:
- Need to understand new concepts
