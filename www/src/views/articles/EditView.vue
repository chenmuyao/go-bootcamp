<template>
  <div class="edit-article bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="edit-article-page w-full max-w-4xl p-6 bg-white rounded-lg shadow-md">
      <div class="flex justify-between items-center mb-6">
        <a
          href="/articles/list"
          @click.prevent="checkUnsavedChanges"
          class="text-gumbo-700 underline hover:text-gumbo-500"
        >
          Back to List
        </a>
        <h1 class="text-3xl font-bold">Edit Article</h1>
      </div>

      <form @submit.prevent="saveArticle">
        <div class="mb-4">
          <label for="title" class="block text-sm font-medium text-gray-700">Title</label>
          <input
            id="title"
            v-model="article.title"
            type="text"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-gumbo-500 focus:border-gumbo-500"
          />
        </div>

        <div class="mb-4">
          <label for="content" class="block text-sm font-medium text-gray-700">Content</label>
          <textarea
            id="content"
            v-model="article.content"
            rows="6"
            class="mt-1 block w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:ring-gumbo-500 focus:border-gumbo-500"
          ></textarea>
        </div>

        <div class="flex space-x-4">
          <button
            type="submit"
            class="px-4 py-2 bg-gumbo-500 text-white rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50"
          >
            Save
          </button>
          <button
            type="button"
            @click="publishArticle"
            class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-400"
          >
            Publish
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import axios from "@/axios/axios";
import type { Article } from '@/types/article';

export default defineComponent({
  name: 'EditView',
  setup() {
    const article = ref<Article>({
      id: 0,
      title: '',
      content: '',
      likeCnt: 0,
      liked: false,
      collectCnt: 0,
      collected: false,
      readCnt: 0,
    });

    const unsavedChanges = ref(false);

    const saveArticle = async () => {
      try {
        await axios.post('/articles/edit', article.value);
        unsavedChanges.value = false;
        alert('Article saved successfully!');
      } catch (error) {
        console.error('Error saving article:', error);
      }
    };

    const publishArticle = async () => {
      try {
        await axios.post('/articles/publish', article.value);
        unsavedChanges.value = false;
        alert('Article published successfully!');
      } catch (error) {
        console.error('Error publishing article:', error);
      }
    };

    const checkUnsavedChanges = () => {
      if (unsavedChanges.value) {
        if (confirm('You have unsaved changes. Do you want to leave without saving?')) {
          window.location.href = '/articles/list';
        }
      } else {
        window.location.href = '/articles/list';
      }
    };

    const loadArticle = async (id: number) => {
      try {
        const response = await axios.get(`/articles/detail/${id}`);
        article.value = response.data.data;
      } catch (error) {
        console.error('Error loading article:', error);
      }
    };

    onMounted(() => {
      const params = new URLSearchParams(window.location.search);
      const id = params.get('id');
      if (id) {
        loadArticle(Number(id));
      }
    });

    return {
      article,
      saveArticle,
      publishArticle,
      checkUnsavedChanges,
    };
  },
});
</script>

<style scoped>
.edit-article {
  background-color: #f3f4f6;
}
</style>
