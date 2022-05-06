<template>
  <slot></slot>
</template>

<script lang="ts">
import { useMedia } from '@/composables/media';
import { defineComponent, watchEffect } from 'vue';

export default defineComponent({
  name: 'ThemeProvider',
  setup() {
    const darkMode = useMedia('(prefers-color-scheme: dark)');

    watchEffect(() => {
      if (
        localStorage.theme === 'dark' ||
        (!('theme' in localStorage) && darkMode.value)
      ) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
    });

    return {};
  },
});
</script>
