<template>
  <div class="flex flex-col bg-pbg-500 w-full h-screen">
    <h1 class="text-center text-ptextd-500 text-4xl m-3 font-bold">🤖 Product Photography, but Crowdsourced ⚡️</h1>
    <div class="flex flex-row w-full h-full">
      <div class="flex flex-col py-4 pl-4 bg-pbg-500 basis-7/12">
        <div class="flex flex-row basis-5/12">
          <!-- product image -->
          <div class="bg-pbrown-500 p-3 mr-2 basis-1/2 flex flex-col rounded-3xl shadow-2xl shadow-red-500">
            <input class="text-ptxtl-500 h-[45px] text-sm" type="file" @change="handleProductFileChange" accept=".png" />
            <input class="rounded-xl h-[35px] p-2" v-model="productPrompt" placeholder="Enter Prompt" />
            <div v-if="!productImageUrl" class="mt-2 flex items-center h-full w-full flex-col justify-center">
              <h1 class="text-3xl text-ptxtl-500 text-center w-full font-bold">1. Upload Product Image &</h1>
              <h1 class="text-3xl text-ptxtl-500 text-center w-full font-bold">Enter Prompt</h1>
            </div>
            <div class="flex flex-col items-center" v-if="productImageUrl">
              <img class="mt-2 h-[160px] w-auto" :src="productImageUrl" alt="Uploaded Image" />
            </div>
          </div>
          <!-- logo image -->
          <div class="bg-pbrown2-500 p-3 ml-2 basis-1/2 flex flex-col rounded-3xl shadow-2xl shadow-yellow-500">
            <input class="text-ptxtl-500 h-[45px] text-sm" type="file" @change="handleLogoFileChange" accept=".png" />
            <input class="rounded-xl h-[35px] p-2" v-model="logoPrompt" placeholder="Enter Prompt" />
            <div v-if="!logoImageUrl" class="flex items-center h-full w-full flex-col justify-center">
              <h1 class="text-3xl text-ptxtl-500 text-center w-full font-bold">2. Upload Logo Image &</h1>
              <h1 class="text-3xl text-ptxtl-500 text-center w-full font-bold">Enter Prompt</h1>
            </div>
            <div class="flex flex-col items-center" v-if="logoImageUrl">
              <img class="mt-2 h-[160px]" :src="logoImageUrl" alt="Uploaded Image" />
            </div>
          </div>
        </div>
        <!-- task list -->
        <div class="bg-pbrown3-500 mt-4 basis-7/12 overflow-auto flex flex-col rounded-3xl shadow-2xl shadow-blue-500">
          <div v-if="!productImageUrl || !logoImageUrl || productPrompt.length === 0 || logoPrompt.length === 0" class="flex items-center h-full w-full">
            <h1 class="text-3xl text-ptxtl-500 text-center w-full font-bold">3. Submit Job Requests</h1>
          </div>
          <div class="flex items-center justify-center h-full w-full" v-if="productImageUrl && logoImageUrl && productPrompt.length > 0 && logoPrompt.length > 0 && !jobPending">
            <button @click="submitJob" class="bg-btna-500 text-btntxta-500 w-[100px] h-[50px] rounded-xl font-bold">Submit 🚀</button>
          </div>
          <div v-if="jobPending" class="flex flex-col h-full w-full p-3">
            <h1 class="text-ptxtl-500 text-center text-3xl font-bold mb-2">{{ pendingTasks }} Tasks Pending</h1>
            <div class="flex flex-col w-full h-full">
              <DVMTask :title="'1. Remove Product Background'" :jobRequest="jobClearBackground"></DVMTask>
              <DVMTask :title="'2. Generate Logo'" :jobRequest="jobGenerateLogo"></DVMTask>
              <DVMTask :title="'3. Generate Background'" :jobRequest="jobGenerateBackground"></DVMTask>
              <DVMTask :title="'4. Overlay Product on Background'" :jobRequest="jobOverlayProduct"></DVMTask>
              <DVMTask :title="'5. Overlay Logo on Product'" :jobRequest="jobOverlayLogo"></DVMTask>
            </div>
          </div>
        </div>
      </div>
      <div class="flex flex-col p-4 basis-5/12">
        <!-- final image -->
        <div class="bg-pbrown4-500 h-full flex flex-col rounded-3xl shadow-green-500 shadow-2xl">
          <div v-if="jobOverlayLogo.output.length === 0" class="flex items-center h-full w-full">
            <h1 class="text-3xl text-ptxtl-500 text-center w-full font-bold">4. View Final Image here!</h1>
          </div>
          <div class="h-full flex flex-col items-center justify-center" v-if="jobOverlayLogo.output.length > 0">
            <img class="h-[600px]" :src="jobOverlayLogo.output" alt="Uploaded Image" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import DVMTask from '../components/DVMTask.vue'
import { type Ref, ref, computed, onMounted } from 'vue'
import { useJobRequest } from '../composables/job_request'
import { useGlobals } from '@/main'

const productImageUrl = ref<string | null>(null);
const productPrompt = ref<string>('')
const logoImageUrl = ref<string | null>(null);
const logoPrompt = ref<string>('')
const jobPending = ref<boolean>(false);
const { 
  submitJobRequests,
  jobClearBackground,
  jobGenerateBackground,
  jobGenerateLogo,
  jobOverlayProduct,
  jobOverlayLogo,
} = useJobRequest()
const { imgbbSecret } = useGlobals()

const handleProductFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];

  if (file && file.type === 'image/png') {
    uploadImage(productImageUrl, file);
  } else {
    // Handle invalid file type or no file selected
  }
};

const handleLogoFileChange = (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];

  if (file && file.type === 'image/png') {
    uploadImage(logoImageUrl, file);
  } else {
    // Handle invalid file type or no file selected
  }
};

const uploadImage = async (imageUrl:Ref<string | null>, file: File) => {
  const apiUrl = 'https://api.imgbb.com/1/upload';
  const apiKey = imgbbSecret; // Replace with your actual API key
  const formData = new FormData();
  formData.append('image', file);

  try {
    const response = await fetch(`${apiUrl}?expiration=600&key=${apiKey}`, {
      method: 'POST',
      body: formData,
    });

    if (response.ok) {
      const responseData = await response.json();
      if (responseData.data && responseData.data.url) {
        imageUrl.value = responseData.data.url;
        console.log(imageUrl.value)
      } else {
        // Handle API response without the expected URL
      }
    } else {
      // Handle non-200 HTTP response status
    }
  } catch (error) {
    // Handle errors in API request
    console.error('Error uploading image:', error);
  }
}

const pendingTasks = computed(() => {
  return jobClearBackground.value.pending +
    jobGenerateBackground.value.pending +
    jobGenerateLogo.value.pending +
    jobOverlayLogo.value.pending +
    jobOverlayProduct.value.pending
})

const submitJob = () => {
  jobPending.value = true
  submitJobRequests(
    productImageUrl.value!,
    productPrompt.value,
    logoImageUrl.value!,
    logoPrompt.value,
  )
}

onMounted(() => {
    checkWebLN()
})

function checkWebLN() {
  return new Promise((resolve) => {
    detectWebLNProvider().then(async (webln: any) => {
      if (webln != null) {
        try {
          // webln available, mark it in user profile
          // ProfileManager.user().value.canWebLN = true
        } catch {
          console.log('failed to enable webln')
        }
      } else {
        console.log('webln not detected')
      }
      resolve(null)
    })
  })
}

async function detectWebLNProvider(timeoutParam: number = 3000) {
  const timeout = timeoutParam
  const interval = 100;
  let handled = false;

  return new Promise((resolve) => {
    if ((window as any).webln) {
      handleWebLN();
    } else {
      document.addEventListener("webln:ready", handleWebLN, { once: true });
      
      let i = 0;
      const checkInterval = setInterval(function() {
        if ((window as any).webln || i >= timeout/interval) {
          handleWebLN();
          clearInterval(checkInterval);
        }
        i++;
      }, interval);
    }

    function handleWebLN() {
      if (handled) {
        return;
      }
      handled = true;

      document.removeEventListener("webln:ready", handleWebLN);

      if ((window as any).webln) {
        resolve((window as any).webln);
      } else {
        resolve(null);
      }
    }
  });
}

</script>
