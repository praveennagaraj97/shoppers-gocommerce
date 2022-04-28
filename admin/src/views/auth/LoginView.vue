<template>
  <div class="h-screen bg-mirage flex items-center justify-center sm:p-0 p-4">
    <div class="relative sm:w-96 w-full">
      <div class="md:block hidden">
        <div
          class="bg-blue-300/5 h-40 w-40 rounded-xl absolute -top-20 -left-20"
        ></div>
        <div
          class="border-blue-300/10 border h-20 w-20 rounded-xl absolute -top-14 left-10"
        ></div>
      </div>
      <div class="md:block hidden">
        <div
          class="bg-blue-300/5 h-40 w-40 rounded-xl absolute -bottom-20 -right-20"
        ></div>
        <div
          class="border-blue-300/10 border h-20 w-20 rounded-xl absolute -bottom-14 right-10"
        ></div>
      </div>

      <div class="w-full bg-slate-600 rounded-lg p-2 drop-shadow-2xl">
        <img
          src="../../assets/brand_logo.png"
          alt=""
          width="180"
          class="mx-auto"
        />
        <div class="mt-4 px-2 text-left">
          <h1 class="text-lg font-semibold text-gray-200 mb-2">
            Welcome to {{ brandName }}
          </h1>
          <p class="text-gray-300 mb-3">
            Please sign-in to your account and start managing
          </p>

          <div class="text-gray-100 mb-3">
            <label for="email" class="text-sm">Email</label>
            <input
              type="text"
              class="w-full p-2 rounded-lg shadow-xl input-focus"
              placeholder="Enter your email"
            />
          </div>

          <div class="text-gray-100 mb-2">
            <label for="password" class="text-sm">Password</label>
            <input
              type="password"
              class="w-full p-2 rounded-lg shadow-xl input-focus"
              placeholder="Enter your password"
            />
          </div>

          <div class="text-gray-100 flex space-x-1 items-center">
            <input type="checkbox" />
            <label for="remember_me" class="text-sm">Remember Me</label>
          </div>

          <div class="my-4">
            <p class="text-gray-300 text-sm text-center">
              By clicking on below you agree to our
              <a class="text-mandy"> terms and conditions </a>
            </p>
          </div>

          <button
            class="will-change-transform py-2 w-full text-gray-100 mt-4 mb-2 bg-cloud_burst rounded-md"
          >
            Login
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
  setup() {
    const brandName = process.env.VUE_APP_BRAND_NAME || 'Brand Name';
    return { brandName };
  },
});
</script>
