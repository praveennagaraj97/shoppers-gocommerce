<template>
  <ThemeProvider>
    <router-view></router-view>
  </ThemeProvider>
</template>

<script lang="ts">
import ThemeProvider from '@/components/providers/ThemeProvider.vue';
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'App',
  components: { ThemeProvider },
});
</script>
