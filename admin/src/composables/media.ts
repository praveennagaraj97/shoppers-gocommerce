import { onBeforeMount, onUnmounted, ref } from 'vue';

export const useMedia = (query: string) => {
  const matches = ref(true);
  const media = window.matchMedia(query);

  const onChange = () => {
    console.log('changed');
    matches.value = media.matches;
  };

  onBeforeMount(() => {
    if (media.matches !== matches.value) {
      matches.value = media.matches;
    }

    media.addEventListener('change', onChange);
  });

  onUnmounted(() => media.removeEventListener('change', onChange));

  return matches;
};
