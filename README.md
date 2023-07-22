# PhotoBolt
Your magical product photographer

## Instructions
- Install photobolt fork of `automatic1111` (https://github.com/lnconsole/stable-diffusion-webui)
  - Install photon and dreamshaper checkpoints
  - Install ColoredIcon Lora
  - Enable ControlNet in the UI, then download Canny preprocessor
  - git checkout `photobolt` branch
  - Launch with `./webui.sh --nowebui`
- Install `rembg` via pip (https://github.com/danielgatis/rembg)
- Install `ffmpeg` (brew or https://ffmpeg.org/download.html)
- `cp .env-example .env` in the `/env` folder. Populate it. NOTE: automatic1111 url should be pointing to the api server (7681 by default)
- `go run .`

## TODO
- Implement service components (automatic1111)
- gin server
  - POST /background (generate a new background for an image based on provided prompt)
  - POST /icon (generate a new icon w/ or w/e a provided image)
  - POST /overlay (combine 2 images together, based on provided position)
- hook it up to aperture??
