<template>
  <v-container class="d-flex justify-center align-center h-100">
    <v-card class="w-50 h-auto">
      <v-card-item>
        <v-card-title class="text-h5">
          Login
        </v-card-title>
      </v-card-item>
      <v-card-text class="pt-2">
        <form @submit.prevent="login" fast-fail>
          <v-row>
            <v-col cols="12">
              <v-text-field v-model="username" label="Username" outlined dense required
              ></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-text-field v-model="password" label="Password" outlined dense required type="password"></v-text-field>
              <v-btn color="primary" type="submit">
                Login
              </v-btn>
            </v-col>
            <v-col cols="12" class="text-center text-base">
              <span>Do not have account? </span>
              <router-link to="/register">Register</router-link>
            </v-col>
          </v-row>
        </form>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAppStore } from '@/stores/app';

const router = useRouter();
const appStore = useAppStore();
const username = ref('');
const password = ref('');

const login = async () => {
  await appStore.login(username.value, password.value);
  router.push('/');
}

</script>

<route>
  {
    meta: {
      layout: 'switcher',
    }
  }
</route>