# PhotoBolt
PhotoBolt converts a Product Image into a Poster.

It uses
1. `Generative AI`(Stable Diffusion) for image processing,
2. `Decentralized Communication Protocol`(nostr) for task outsourcing,
3. `Bitcoin Technology`(Lightning Network) for payment processing

<img width="1494" alt="Screenshot 2023-07-29 at 1 08 31 PM" src="https://github.com/lnconsole/PhotoBolt/assets/43709958/b5673eef-9526-44d9-a4c0-b51d8c223f31">
.
  
This repo contains both the `client`(vue) and `service provider`(go). Client breaks down the poster generation task into 5 smaller tasks, chain them together via [NIP90(Data Vending Machine)](https://github.com/nostr-protocol/nips/blob/vending-machine/90.md) and broadcast them to the nostr network. Service provider accepts each job requests by prompting for a payment, then process and return the job result back to the client.

Although each job request may depend on another job request as an input, they could still be processed independently by different service providers. You could test it out by running two PhotoBolt service provider instances when generating a poster. You should be able to tell which service provider took a task based on the avatar rendered under the `Tasks Pending` UI section.

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
- You'll also need an [ImgBB](https://imgbb.com/) api secret key
- `cp .env-example .env` in the `/env` folder. Populate it. NOTE: automatic1111 url should be pointing to the api server (7681 by default)
- `go run .` to start the service provider
- `cp .env-example .env` in the `/frontend` folder. Populate it
- `cd frontend; npm run dev` to start the client

## L402
This project started with L402 and eventually moved towards NIP90. To test L402, get `aperture` running, and then test out the CLI client by running `cd client; go run .`
