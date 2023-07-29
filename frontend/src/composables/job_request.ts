import { ref, type Ref, computed } from 'vue'
import { relayInit, generatePrivateKey, getPublicKey, getEventHash, getSignature, Kind } from 'nostr-tools'

export enum JobRequestState {
    Preparing,
    Submitted,
    Offered,
    Processing,
    Completed,
}

export interface JobRequest {
    eventID: string
    state: JobRequestState
    stateVerbose: string
    invoice: string
    output: string
    pending: number
    pk: string
}

const jobClearBackground: Ref<JobRequest> = ref({
    eventID: '',
    state: JobRequestState.Preparing,
    stateVerbose: 'Sending Job Request...',
    invoice: '',
    output: '',
    pending: 1,
    pk: '',
})
const jobGenerateLogo: Ref<JobRequest> = ref({
    eventID: '',
    state: JobRequestState.Preparing,
    stateVerbose: 'Sending Job Request...',
    invoice: '',
    output: '',
    pending: 1,
    pk: '',
})
const jobGenerateBackground: Ref<JobRequest> = ref({
    eventID: '',
    state: JobRequestState.Preparing,
    stateVerbose: 'Sending Job Request...',
    invoice: '',
    output: '',
    pending: 1,
    pk: '',
})
const jobOverlayProduct: Ref<JobRequest> = ref({
    eventID: '',
    state: JobRequestState.Preparing,
    stateVerbose: 'Sending Job Request...',
    invoice: '',
    output: '',
    pending: 1,
    pk: '',
})
const jobOverlayLogo: Ref<JobRequest> = ref({
    eventID: '',
    state: JobRequestState.Preparing,
    stateVerbose: 'Sending Job Request...',
    invoice: '',
    output: '',
    pending: 1,
    pk: '',
})

const sk = generatePrivateKey() // `sk` is a hex string
const pk = getPublicKey(sk) // `pk` is a hex string
let jobreq_state = new Map<string, Ref<JobRequest>>();

export function useJobRequest() {
    const submitJobRequests = async (
        productUrl: string, 
        productPrompt: string, 
        logoUrl: string,
        logoPrompt: string, 
    ) => {
        const relay = relayInit('ws://127.0.0.1:7447')
        relay.on('connect', () => {
          console.log(`connected to ${relay.url}`)
        })
        relay.on('error', () => {
          console.log(`failed to connect to ${relay.url}`)
        })
        await relay.connect()

        // clear background
        const clearBackgroundEvent = imageManipulationEvent([
            ["i", productUrl, "url"],
            ["param", "background", "clear"],
        ])
        jobClearBackground.value.eventID = clearBackgroundEvent.id
        jobreq_state.set(clearBackgroundEvent.id, jobClearBackground)

        // generate logo
        const generateLogoEvent = imageGenerationEvent([
            ["i", logoUrl, "url"],
        ], logoPrompt)
        jobGenerateLogo.value.eventID = generateLogoEvent.id
        jobreq_state.set(generateLogoEvent.id, jobGenerateLogo)

        // generate background
        const generateBackgroundEvent = imageGenerationEvent([
            ["i", clearBackgroundEvent.id, "job"],
            ["param", "control-net", "canny"],
        ], productPrompt)
        jobGenerateBackground.value.eventID = generateBackgroundEvent.id
        jobreq_state.set(generateBackgroundEvent.id, jobGenerateBackground)

        // overlay product
        const overlayProductEvent = imageManipulationEvent([
            ["i", clearBackgroundEvent.id, "job", "front"],
            ["i", generateBackgroundEvent.id, "job", "back"],
            ["param", "overlay", "full"],
        ])
        jobOverlayProduct.value.eventID = overlayProductEvent.id
        jobreq_state.set(overlayProductEvent.id, jobOverlayProduct)

        // overlay logo
        const overlayLogoEvent = imageManipulationEvent([
            ["i", generateLogoEvent.id, "job", "front"],
            ["i", overlayProductEvent.id, "job", "back"],
            ["param", "overlay", "logo"],
        ])
        jobOverlayLogo.value.eventID = overlayLogoEvent.id
        jobreq_state.set(overlayLogoEvent.id, jobOverlayLogo)

        // subscribe
        let sub = relay.sub([
            {
                since: Math.floor(Date.now() / 1000),
                kinds: [65000, 65001],
            }
        ])
        sub.on('event', event => {
            console.log('got an event:', event)
            const pTag = event.tags.find(([k]) => k === 'p')
            if (pTag === undefined || pTag[1] !== pk) {
                console.log('what is this: ' + pTag)
                return
            }
            const eTag = event.tags.find(([k]) => k === 'e')
            if (eTag === undefined || !jobreq_state.has(eTag[1])) {
                console.log('what is this: ' + eTag)
                return
            }
            let state = jobreq_state.get(eTag[1])!
            state.value.pk = event.pubkey
            if (event.kind === 65000) {
                const statusTag = event.tags.find(([k]) => k === 'status')
                if (statusTag === undefined) {
                    console.log('what is this status: ' + statusTag)
                    return
                }
                state.value.stateVerbose = statusTag[2]
                if (statusTag[1] === 'payment-required') {
                    const amountTag = event.tags.find(([k]) => k === 'amount')!
                    state.value.state = JobRequestState.Offered
                    state.value.invoice = amountTag[2]
                } else if (statusTag[1] === 'processing') {
                    state.value.state = JobRequestState.Processing
                } else {
                    console.log('what is this status: ' + statusTag)
                    return
                }
            } else if (event.kind === 65001) {
                state.value.state = JobRequestState.Completed
                state.value.stateVerbose = 'Job Completed!'
                state.value.output = event.content
                state.value.pending = 0
            } else {
                console.log('what is this kind: ' + event.kind)
                return
            }
        })
        

        // send job requests
        let clearBackgroundPub = relay.publish(clearBackgroundEvent)
        clearBackgroundPub.on('ok', () => {
            console.log(`${relay.url} has accepted our event`)
        })
        clearBackgroundPub.on('failed', (reason: any) => {
            console.log(`failed to publish to ${relay.url}: ${reason}`)
        })
        
        let generateLogoPub = relay.publish(generateLogoEvent)
        generateLogoPub.on('ok', () => {
            console.log(`${relay.url} has accepted our event`)
        })
        generateLogoPub.on('failed', (reason: any) => {
            console.log(`failed to publish to ${relay.url}: ${reason}`)
        })

        let generateBackgroundPub = relay.publish(generateBackgroundEvent)
        generateBackgroundPub.on('ok', () => {
            console.log(`${relay.url} has accepted our event`)
        })
        generateBackgroundPub.on('failed', (reason: any) => {
            console.log(`failed to publish to ${relay.url}: ${reason}`)
        })

        let overlayProductPub = relay.publish(overlayProductEvent)
        overlayProductPub.on('ok', () => {
            console.log(`${relay.url} has accepted our event`)
        })
        overlayProductPub.on('failed', (reason: any) => {
            console.log(`failed to publish to ${relay.url}: ${reason}`)
        })

        let overlayLogoPub = relay.publish(overlayLogoEvent)
        overlayLogoPub.on('ok', () => {
            console.log(`${relay.url} has accepted our event`)
        })
        overlayLogoPub.on('failed', (reason: any) => {
            console.log(`failed to publish to ${relay.url}: ${reason}`)
        })

        jobClearBackground.value.state = JobRequestState.Submitted
        jobClearBackground.value.stateVerbose = 'Job Request Sent! Waiting for Service Providers...'
        jobGenerateBackground.value.state = JobRequestState.Submitted
        jobGenerateBackground.value.stateVerbose = 'Job Request Sent! Waiting for Service Providers...'
        jobGenerateLogo.value.state = JobRequestState.Submitted
        jobGenerateLogo.value.stateVerbose = 'Job Request Sent! Waiting for Service Providers...'
        jobOverlayProduct.value.state = JobRequestState.Submitted
        jobOverlayProduct.value.stateVerbose = 'Job Request Sent! Waiting for Service Providers...'
        jobOverlayLogo.value.state = JobRequestState.Submitted
        jobOverlayLogo.value.stateVerbose = 'Job Request Sent! Waiting for Service Providers...'
    }

    const imageGenerationEvent = (tags: string[][], content: string): any => {
        let event: any = {
            kind: 65005,
            created_at: Math.floor(Date.now() / 1000),
            tags: tags,
            content: content,
            pubkey: pk,
        }

        event.id = getEventHash(event)
        event.sig = getSignature(event, sk)

        return event
    }

    const imageManipulationEvent = (tags: string[][]): any => {
        let event: any = {
            kind: 65007,
            created_at: Math.floor(Date.now() / 1000),
            tags: tags,
            content: '',
            pubkey: pk,
        }

        event.id = getEventHash(event)
        event.sig = getSignature(event, sk)

        return event
    }

    return {
        submitJobRequests,
        jobClearBackground,
        jobGenerateLogo,
        jobGenerateBackground,
        jobOverlayProduct,
        jobOverlayLogo,
        JobRequestState,
    }
}