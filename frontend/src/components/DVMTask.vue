<template>
    <div class="bg-pbg2-500 basis-1/5 mx-1 mb-3 flex flex-row rounded-xl">
        <div class="flex flex-col basis-10/12">
            <div class="basis-1/2 h-full ml-1 mt-1 font-bold">
                {{ props.title }}
            </div>
            <div class="basis-1/2 h-full ml-1 mb-1 flex flex-row items-center justify-left">
                <img v-if="props.jobRequest.pk.length > 0" class="w-5 h-5 my-1 ml-2 rounded-full" :src="`https://robohash.org/${props.jobRequest.pk}.png?set=set5`" />
                <p v-if="props.jobRequest.pk.length > 0">:</p>
                <p class="ml-3">{{ props.jobRequest.stateVerbose }}</p>
            </div>
        </div>
        <div class="flex flex-col justify-center items-center basis-2/12">
            <div v-show="props.jobRequest.state === JobRequestState.Offered">
                <button @click="pay" class="bg-btna-500 text-btntxta-500 font-bold w-[80px] h-[40px] rounded-xl">Pay</button>
            </div>
            <div v-show="props.jobRequest.state === JobRequestState.Completed">
                <img class="h-[60px]" :src="props.jobRequest.output"/>
            </div>
            <div v-show="props.jobRequest.state === JobRequestState.Processing">
                <p class="text-3xl">⏱️</p>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { type JobRequest, JobRequestState } from '@/composables/job_request';
import {type PropType} from 'vue'

const props = defineProps({
    jobRequest: {
        type: Object as PropType<JobRequest>,
        required: true,
    },
    title: {
        type: String,
        required: true,
    },
});

const pay = async () => {
    const inv = props.jobRequest.invoice
    console.log("hi")
    try {
        await (window as any).webln.enable()
        await (window as any).webln.sendPayment(inv)
    } catch (err) {
        console.log('failed to pay bolt11:' + inv)
    }
    console.log("h2")
}

</script>