<template>
  <div class="article-read bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="article-page w-full max-w-5xl p-6 bg-white rounded-lg shadow-md">
      <!-- Back to List Link -->
      <div class="mb-6">
        <a
          href="/articles/list"
          class="text-gumbo-500 hover:text-gumbo-700 inline-flex items-center"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5 mr-2"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z"
              clip-rule="evenodd"
            />
          </svg>
          Back to Articles
        </a>
      </div>

      <!-- Article Content -->
      <article v-if="article" class="prose max-w-none">
        <h1 class="text-3xl font-bold text-gumbo-700 mb-4">{{ article.title }}</h1>

        <div
          class="article-content prose prose-gumbo max-w-none"
          v-html="article.content"
        ></div>

        <!-- Article Interaction Icons -->
        <div class="flex justify-between items-center mt-6 border-t pt-4 border-gray-200">
          <div class="flex space-x-4">
            <!-- Read Count -->
            <div class="flex items-center text-gumbo-500">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 mr-2"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                <path
                  fill-rule="evenodd"
                  d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z"
                  clip-rule="evenodd"
                />
              </svg>
              <span>{{ article.readCnt }}</span>
            </div>

            <!-- Like Count -->
            <div
              class="flex items-center cursor-pointer text-gumbo-500 hover:text-gumbo-700"
              @click="toggleLike"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 mr-2"
                :class="{ 'text-red-500': article.liked }"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path
                  fill-rule="evenodd"
                  d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z"
                  clip-rule="evenodd"
                />
              </svg>
              <span>{{ article.likeCnt }}</span>
            </div>

            <!-- Collect Count -->
            <div
              class="flex items-center cursor-pointer text-gumbo-500 hover:text-gumbo-700"
              @click="toggleCollect"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 mr-2"
                :class="{ 'text-gumbo-700': article.collected }"
                viewBox="0 0 20 20"
                fill="currentColor"
              >
                <path d="M5 4a2 2 0 012-2h6a2 2 0 012 2v14l-5-2.5L5 18V4z" />
              </svg>
              <span>{{ article.collectCnt }}</span>
            </div>
          </div>
        </div>
      </article>

      <!-- Loading or Not Found State -->
      <div v-else class="text-center text-gumbo-500">
        Loading article...
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import axios from '@/axios/axios';
import type { Article } from '@/types/article';

export default defineComponent({
  name: 'ReadView',
  setup() {
    const article = ref<Article | null>(null);

    const fetchArticle = async () => {
      try {
        const params = new URLSearchParams(window.location.search);
        const id = params.get('id');
        const response = await axios.get(`/articles/pub/${id}`);
        article.value = response.data.data;
      } catch (error) {
        console.error('Error fetching article:', error);
      }
    };

    const toggleLike = async () => {
      if (!article.value) return;

      try {
        const liked = !article.value.liked;
        const response = await axios.post('/articles/pub/like', {
          id: article.value.id,
          liked
        });

        if (response.status === 200) {
          article.value.liked = liked;
          article.value.likeCnt += liked ? 1 : -1;
        }
      } catch (error) {
        console.error('Error toggling like:', error);
      }
    };

    const toggleCollect = async () => {
      if (!article.value) return;

      try {
        const collected = !article.value.collected;
        const response = await axios.post('/articles/pub/collect', {
          id: article.value.id,
          collected
        });

        if (response.status === 200) {
          article.value.collected = collected;
          article.value.collectCnt += collected ? 1 : -1;
        }
      } catch (error) {
        console.error('Error toggling collect:', error);
      }
    };

    onMounted(() => {
      fetchArticle();
    });

    return {
      article,
      toggleLike,
      toggleCollect
    };
  }
});
</script>

<style scoped>
.articles {
  background-color: #f9fafb;
}

.prose-gumbo {
  @apply text-gray-800 leading-relaxed;
}

.prose-gumbo h1,
.prose-gumbo h2,
.prose-gumbo h3 {
  @apply text-gumbo-700 font-bold;
}

.prose-gumbo p {
  @apply mb-4;
}
</style>
