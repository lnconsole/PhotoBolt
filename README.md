# PhotoBolt
PhotoBolt converts a Product Image into a Poster

This repo contains both the `client`(vue) and `service provider`(go). Client breaks down the poster generation task into 5 smaller tasks, chain them together via [NIP90(Data Vending Machine)](https://github.com/nostr-protocol/nips/blob/vending-machine/90.md) and broadcast them to the nostr network. Service provider accepts each job requests by prompting for a payment, then process and return the job result back to the client.

Although each job request may depend on another job request as an input, they could still be processed independently by different service providers. You could test it out by running two PhotoBolt service provider instances when generating a poster. You should be able to tell which service provider took a task based on the avatar rendered under the `Pending Tasks` UI section.

Video Demo [here]()

## Summary of tech used
- [automatic1111](https://github.com/AUTOMATIC1111/stable-diffusion-webui) Stable Diffusion server
- [rembg](https://github.com/danielgatis/rembg) Background removal tool
- [ffmpeg](https://ffmpeg.org/) Media processing tool
- [polar](https://github.com/jamaljsr/polar) Simulated Lightning Network
- [alby](https://getalby.com/#alby-extension) webLN browser wallet integration
- [relayer](https://github.com/fiatjaf/relayer) nostr relay
- [NIP90](https://github.com/nostr-protocol/nips/blob/vending-machine/90.md) nostr Data Vending Machine Proposal
- [aperture](https://github.com/lightninglabs/aperture) L402 gateway server

## Instructions
- Install `automatic1111`
  - Install photon and dreamshaper checkpoints
  - Install ColoredIcon Lora
  - Enable ControlNet in the UI, then download Canny preprocessor
  - git checkout `photobolt` branch
  - Launch with `./webui.sh --nowebui`
- Install `rembg`
- Install `ffmpeg`
- Install `polar` or a LND mainnet node. Both client and server use Lightning
- Install `alby` (you will need a mainnet node for the server)
- Connect to a public nostr relay or Install `relayer`
- `cp .env-example .env` in the `/env` folder. Populate it. NOTE: automatic1111 url should be pointing to the api server (7681 by default)
- `go run .` to start the service provider
- `cd frontend; npm run dev` to start the client

## L402
This project started with L402 and eventually moved towards NIP90. To test L402, get `aperture` running, and then test out the CLI client by running `cd client; go run .`
