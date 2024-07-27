<template>
  <v-container class="d-flex justify-center align-center h-100">
    <v-card class="w-50 h-auto">
      <v-card-item>
        <v-card-title class="text-h5">
          Register
        </v-card-title>
      </v-card-item>
      <v-card-text class="pt-2">
        <form @submit.prevent="register">
          <v-row>
            <v-col cols="12">
              <v-text-field v-model="username" label="Username" outlined dense required></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-text-field v-model="password" label="Password" outlined dense required type="password"></v-text-field>
              <v-btn color="primary" type="submit">
                Register
              </v-btn>
            </v-col>
          </v-row>
        </form>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAppStore} from '@/stores/app';

const router = useRouter();
const appStore = useAppStore();
const username = ref('');
const password = ref('');

const register = async () => {
  await appStore.register(username.value, password.value);
  router.push('/login');
}
</script>

<route>
  {
    meta: {
      layout: 'switcher',
    }
  }
</route>