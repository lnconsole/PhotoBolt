# PhotoBolt
PhotoBolt converts a Product Image into a Poster.

It uses
1. **Generative AI** *(Stable Diffusion)* for `image processing`,
2. **Decentralized Communication Protocol** *(nostr)* for `task planning & outsourcing`,
3. **Bitcoin Technology** *(Lightning Network)* for `payment processing`

<img width="1494" alt="Screenshot 2023-07-29 at 1 08 31 PM" src="https://github.com/lnconsole/PhotoBolt/assets/43709958/8b2a5966-d402-49d9-b95e-dd6550db98e0">
.
  
This repo contains both the `client`(vue) and `service provider`(go). Client breaks down the poster generation task into 5 smaller tasks, chain them together via [NIP90(Data Vending Machine)](https://github.com/nostr-protocol/nips/blob/vending-machine/90.md) and broadcast them to the nostr network. Service provider accepts each job requests by prompting for a payment, then process and return the job result back to the client.

Although each job request may depend on another job request as an input, they could still be processed independently by different service providers. You could test it out by running two PhotoBolt service provider instances when generating a poster. You should be able to tell which service provider took a task based on the avatar rendered under the `Tasks Pending` UI section.

Video Demo [here](https://youtu.be/xex9rEsrU5I)

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
- Connect to a public nostr relay or Install `relayer` and point both client & service provider to the relayer instance
- You'll also need an [ImgBB](https://imgbb.com/) api secret key
- `cp .env-example .env` in the `/env` folder. Populate it. NOTE: automatic1111 url should be pointing to the api server (7681 by default)
- `go run .` to start the service provider
- `cp .env-example .env` in the `/frontend` folder. Populate it
- `cd frontend; npm run dev` to start the client

Once the client is running, simply follow the instructions on the screen. Upload a product image, and provide a simple prompt. Next upload a logo image and provide a simple prompt. Click `Submit` and pay the invoice whenever there is a job offer (You'll need webLN. If not, go to browser console to find the bolt11). Once all the tasks are completed you should see the final Poster image.

## PhotoBolt Job Chain (rough illustration)
<img width="591" alt="Screenshot 2023-07-29 at 6 01 01 PM" src="https://github.com/lnconsole/PhotoBolt/assets/43709958/53795b55-709f-410e-924d-5b0ad0236cd2">

## L402
This project started with L402 and eventually prioritizes NIP90. To test L402, get `aperture` running, and then test out the CLI client by running `cd client; go run .`

Video Demo [here](https://www.youtube.com/watch?v=TsCNUxBWcvg)
