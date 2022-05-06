<template>
  <slot></slot>
  <Teleport to="#theme-provider">
    <Transition name="theme-switch-slide">
      <div v-if="shouldMount" class="theme_provider_box">
        <div
          class="absolute top-0 cursor-pointer right-0 hover:scale-110 smooth"
          @click="shouldMount = false"
        >
          <IoIoClose className="w-6 h-6" />
        </div>
        <h6 className="text-lg">Choose a style</h6>
        <p className="text-sm">Light or Dark?</p>
        <small>Customise your interface</small>
        <div className="flex justify-between items-center px-2 my-2">
          <small>Light</small>
          <div
            class="w-8 h-5 flex items-center bg-gray-300 rounded-full p-1 cursor-pointer"
          >
            <div class="" v-if="darkMode" @click="darkMode = false">
              <MDLightMode className="w-4 h-4 text-red-500" />
            </div>

            <div class="" v-if="!darkMode">
              <MDDarkMode
                className="w-4 h-4 text-red-500"
                @click="darkMode = true"
              />
            </div>
          </div>
          <small>Dark</small>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script lang="ts">
import { useMedia } from '@/composables/media';
import { defineComponent, onMounted, ref, watchEffect } from 'vue';
import IoIoClose from '../icons/io/IoIoClose.vue';
import MDDarkMode from '../icons/md/MDDarkMode.vue';
import MDLightMode from '../icons/md/MDLightMode.vue';

export default defineComponent({
  name: 'ThemeProvider',
  setup() {
    createMountPoint();
    const darkMode = useMedia('(prefers-color-scheme: dark)');
    const shouldMount = ref<boolean>(false);

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

    onMounted(() => {
      shouldMount.value = true;
    });

    return { shouldMount, darkMode };
  },
  components: { IoIoClose, MDLightMode, MDDarkMode },
});

function createMountPoint() {
  const mountPoint = document.createElement('div');
  mountPoint.id = 'theme-provider';
  document.body.appendChild(mountPoint);
}
</script>

<style lang="postcss" scoped>
.theme_provider_box {
  @apply fixed w-52 bottom-3 dark:text-gray-100 px-4 py-2 rounded-lg shadow-2xl z-10 bg-white dark:bg-gray-700 right-2;
}

.theme-switch-slide-enter-active {
  @apply transform duration-700 transition-all opacity-100 translate-x-0 ease-in-out;
}

.theme-switch-slide-leave-active {
  @apply transform  transition-all opacity-0 translate-x-full duration-700;
}

.theme-switch-slide-enter-from,
.theme-switch-slide-leave-to {
  @apply transform  transition-all opacity-0 translate-x-full duration-700;
}
</style>
