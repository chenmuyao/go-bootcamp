<template>
  <div class="articles bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="articles-page w-full max-w-5xl p-6 bg-white rounded-lg shadow-md">
      <div class="flex justify-between items-center mb-6">
        <h1 class="text-3xl font-bold">Articles</h1>
        <button
          @click="goToNewArticle"
          class="px-4 py-2 bg-gumbo-500 text-white rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50">
          New Article
        </button>
      </div>

      <div v-if="articles.length === 0" class="text-center text-gumbo-500">
        No articles to display.
      </div>

      <ul v-else class="divide-y divide-gray-200">
        <li
          v-for="article in articles"
          :key="article.id"
          class="py-4 flex justify-between items-center">
          <div>
            <h2 class="text-xl font-semibold" :class="getStatusClass(article.status)">{{ article.title }}</h2>
            <p class="text-sm text-gray-500">{{ article.abstract }}</p>
          </div>
          <div class="flex items-center space-x-2">
            <div class="flex space-x-2">
              <button
                @click="editArticle(article.id)"
                class="px-3 py-2 text-sm bg-gumbo-500 text-white rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50">
                Edit
              </button>
              <button
                @click="viewArticle(article.id)"
                class="px-3 py-2 text-sm bg-blue-500 text-white rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50">
                View
              </button>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import axios from "@/axios/axios";

type ArticleItem = {
  id: bigint;
  title: string;
  status: number;
  abstract: string;
};

export default defineComponent({
  name: 'ListView',
  setup() {
    const articles = ref<ArticleItem[]>([]);

    const fetchArticles = async (offset = 0, limit = 10) => {
      try {
        const response = await axios.post('/articles/list', { offset, limit });
        articles.value = response.data.data;
      } catch (error) {
        console.error('Error fetching articles:', error);
      }
    };

    const editArticle = (id: bigint) => {
      window.location.href = `/articles/edit?id=${id}`;
    };

    const viewArticle = (id: bigint) => {
      window.location.href = `/articles/view?id=${id}`;
    };

    const goToNewArticle = () => {
      window.location.href = '/articles/edit';
    };

    const getStatusClass = (status: number) => {
      return status === 1
        ? 'text-gumbo-300'
        : status === 2
        ? 'text-gumbo-700'
        : 'text-red-500';
    };

    onMounted(() => {
      fetchArticles();
    });

    return {
      articles,
      editArticle,
      viewArticle,
      goToNewArticle,
      getStatusClass,
    };
  },
});
</script>

<style scoped>
.articles {
  background-color: #f9fafb;
}
</style>

