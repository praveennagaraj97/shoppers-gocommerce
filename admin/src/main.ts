import { createPinia } from 'pinia';
import { createApp } from 'vue';
import App from './App.vue';
import i18n from './i18n';
import './registerServiceWorker';
import router from './router';
import './styles/index.css';

const app = createApp(App);

app.use(i18n);
app.use(createPinia());
app.use(router);
app.mount('#app');
