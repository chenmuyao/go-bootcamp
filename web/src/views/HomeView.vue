<template>
  <div class="home bg-gumbo-100 min-h-screen flex items-center justify-center">
    <div class="home-page w-full max-w-4xl p-6 bg-white rounded-lg shadow-md">
      <h1 class="text-5xl font-bold text-center mb-8">Welcome to WeTravel</h1>
      <p class="text-lg text-center leading-relaxed mb-12">
        Discover, share, and manage your travels with friends. Join the WeTravel community to make every trip unforgettable!
      </p>
      <div class="flex flex-col md:flex-row justify-center space-y-4 md:space-y-0 md:space-x-4">
        <a href="/signup" class="inline-block px-6 py-3 text-lg font-semibold text-white bg-gumbo-500 rounded hover:bg-gumbo-700 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50">
          Sign Up
        </a>
        <a href="/login" class="inline-block px-6 py-3 text-lg font-semibold text-gumbo-700 bg-white border border-gumbo-500 rounded hover:bg-gumbo-100 focus:outline-none focus:ring-2 focus:ring-gumbo-600 focus:ring-opacity-50">
          Log In
        </a>
      </div>
      <div class="mt-16">
        <h2 class="text-3xl font-bold text-center mb-8">Why Choose WeTravel?</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <!-- Feature 1 -->
          <div class="feature text-center">
            <h3 class="text-2xl font-semibold mb-4">Easy Travel Story Sharing</h3>
            <p class="text-base leading-relaxed">
              Share photos, videos, and experiences with friends in real-time. Keep everyone updated on your adventures.
            </p>
          </div>
          <!-- Feature 2 -->
          <div class="feature text-center">
            <h3 class="text-2xl font-semibold mb-4">Expense Management</h3>
            <p class="text-base leading-relaxed">
              Track and split expenses effortlessly among your travel group. Say goodbye to complicated spreadsheets and endless messages.
            </p>
          </div>
          <!-- Feature 3 -->
          <div class="feature text-center">
            <h3 class="text-2xl font-semibold mb-4">Collaborative Planning</h3>
            <p class="text-base leading-relaxed">
              Plan trips together with built-in tools for itinerary creation and coordination. Make sure everyone is on the same page.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Additional custom styles can be added here if necessary */
.feature {
  border: 1px solid #e5e7eb;
  padding: 2rem;
  border-radius: 0.5rem;
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
}

.feature:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
</style>
