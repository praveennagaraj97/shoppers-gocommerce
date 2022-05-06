<template>
  <div
    class="h-screen dark:bg-mirage flex items-center justify-center sm:p-0 p-4"
  >
    <h1>{{ getAuthState }}</h1>
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
        <form @submit.prevent="handleSumit" class="mt-4 px-2 text-left">
          <h1 class="text-lg font-semibold text-gray-200 mb-2">
            {{ $t('login.welcome') }} {{ brandName }}
          </h1>
          <p class="text-gray-300 mb-3">
            {{ $t('login.sign_in') }}
          </p>

          <div class="mb-3">
            <label for="email" class="text-sm text-gray-50">{{
              $t('login.email')
            }}</label>
            <input
              type="text"
              v-model="loginValues.email"
              class="text-black w-full p-2 rounded-lg shadow-xl input-focus"
              :placeholder="$t('login.enter_login_email')"
            />
          </div>

          <div class="mb-2">
            <label for="password" class="text-sm text-gray-50">{{
              $t('login.password')
            }}</label>
            <input
              type="password"
              v-model="loginValues.password"
              class="text-black w-full p-2 rounded-lg shadow-xl input-focus"
              :placeholder="$t('login.enter_login_password')"
            />
          </div>

          <div class="text-gray-100 flex space-x-1 items-center">
            <input type="checkbox" />
            <label for="remember_me" class="text-sm">{{
              $t('login.remember_me')
            }}</label>
          </div>

          <div class="my-4">
            <p class="text-gray-300 text-sm text-center">
              {{ $t('login.agree_statement') }}
              <a class="text-mandy"> {{ $t('login.terms_and_conditions') }} </a>
            </p>
          </div>

          <button
            type="submit"
            class="will-change-transform py-2 w-full text-gray-100 mt-4 mb-2 bg-cloud_burst rounded-md"
          >
            {{ $t('login.login') }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { loginAPI } from '@/services/authService';
import { useAuthStore } from '@/store/auth';
import { defineComponent, reactive, watchEffect } from 'vue';

export default defineComponent({
  name: 'LoginView',
  setup() {
    const { isLogged, updateAuthState, getAuthState } = useAuthStore();

    const brandName = process.env.VUE_APP_BRAND_NAME || 'Brand Name';

    const loginValues = reactive<{ email: string; password: string }>({
      email: '',
      password: '',
    });

    watchEffect(() => {
      console.log(isLogged);
    });

    async function handleSumit() {
      try {
        const { data } = await loginAPI(
          loginValues.email,
          loginValues.password
        );
      } catch (error) {
        console.log(error);
      }
    }

    return { brandName, loginValues, handleSumit, getAuthState };
  },
});
</script>
