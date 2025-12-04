<template>
  <div class="min-h-screen bg-theme-bg flex flex-col items-center justify-center px-4">
    <!-- Unsplash Image -->
    <div class="w-full max-w-4xl mb-8">
      <img 
        :src="imageUrl" 
        alt="Dark space scene" 
        class="w-full h-auto rounded-lg shadow-lg object-cover"
        style="max-height: 70vh;"
        @error="handleImageError"
      />
    </div>

    <!-- Random Quote -->
    <div class="max-w-2xl text-center">
      <blockquote class="text-2xl md:text-3xl font-light text-theme-text italic mb-4">
        "{{ currentQuote.text }}"
      </blockquote>
      <p class="text-theme-textLight text-lg">
        — {{ currentQuote.author }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// Collection of dark space Unsplash images (direct URLs)
const images = [
  'https://images.unsplash.com/photo-1446776653964-20c1d3a81b06?w=1920&h=1080&fit=crop&q=80&auto=format',
  'https://images.unsplash.com/photo-1446776877081-d282a0f896e2?w=1920&h=1080&fit=crop&q=80&auto=format',
  'https://images.unsplash.com/photo-1446776811953-23d5a604b1e5?w=1920&h=1080&fit=crop&q=80&auto=format',
  'https://images.unsplash.com/photo-1502134249126-9f3755a50d78?w=1920&h=1080&fit=crop&q=80&auto=format',
  'https://images.unsplash.com/photo-1462331940025-496dfbfc7564?w=1920&h=1080&fit=crop&q=80&auto=format',
  'https://images.unsplash.com/photo-1506443432602-ac2fcd6f54e0?w=1920&h=1080&fit=crop&q=80&auto=format',
  'https://images.unsplash.com/photo-1516339901601-2e1b62dc0c4a?w=1920&h=1080&fit=crop&q=80&auto=format',
  'https://images.unsplash.com/photo-1502134249126-9f3755a50d78?w=1920&h=1080&fit=crop&q=80&auto=format'
]

// Collection of sad quotes
const quotes = [
  {
    text: "The worst part of holding the memories is not the pain. It's the loneliness of it. Memories need to be shared.",
    author: "Lois Lowry"
  },
  {
    text: "I have of late—but wherefore I know not—lost all my mirth.",
    author: "William Shakespeare"
  },
  {
    text: "The loneliest moment in someone's life is when they are watching their whole world fall apart, and all they can do is stare blankly.",
    author: "F. Scott Fitzgerald"
  },
  {
    text: "We are all broken, that's how the light gets in.",
    author: "Ernest Hemingway"
  },
  {
    text: "The tragedy of life is not that it ends so soon, but that we wait so long to begin it.",
    author: "W.M. Lewis"
  },
  {
    text: "There is no greater sorrow than to recall happiness in times of misery.",
    author: "Dante Alighieri"
  },
  {
    text: "The world breaks everyone, and afterward, many are strong at the broken places.",
    author: "Ernest Hemingway"
  },
  {
    text: "It is a far, far better thing that I do, than I have ever done; it is a far, far better rest that I go to than I have ever known.",
    author: "Charles Dickens"
  },
  {
    text: "The only way out of the labyrinth of suffering is to forgive.",
    author: "John Green"
  },
  {
    text: "We are all in the gutter, but some of us are looking at the stars.",
    author: "Oscar Wilde"
  }
]

const imageUrl = ref(images[0])
const currentQuote = ref(quotes[0])

// Handle image loading errors
const handleImageError = () => {
  // Try a different image if current one fails
  const currentIndex = images.indexOf(imageUrl.value)
  const nextIndex = (currentIndex + 1) % images.length
  imageUrl.value = images[nextIndex]
}

// Select random image and quote on mount
onMounted(() => {
  const randomImageIndex = Math.floor(Math.random() * images.length)
  const randomQuoteIndex = Math.floor(Math.random() * quotes.length)
  imageUrl.value = images[randomImageIndex]
  currentQuote.value = quotes[randomQuoteIndex]
})
</script>
