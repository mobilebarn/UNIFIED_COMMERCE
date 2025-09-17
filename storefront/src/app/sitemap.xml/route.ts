import { MetadataRoute } from 'next';
import { gql } from '@apollo/client';
import { apolloClient } from '@/lib/apollo';

// GraphQL queries for dynamic sitemap generation
const GET_PRODUCTS = gql`
  query GetProducts($limit: Int) {
    products(filter: { limit: $limit }) {
      id
      handle
      updatedAt
    }
  }
`;

const GET_CATEGORIES = gql`
  query GetCategories($limit: Int) {
    categories(filter: { limit: $limit }) {
      id
      handle
      updatedAt
    }
  }
`;

// Base URLs that should always be included
const BASE_URLS: MetadataRoute.Sitemap = [
  {
    url: 'http://localhost:3000',
    lastModified: new Date(),
    changeFrequency: 'yearly',
    priority: 1.0,
  },
  {
    url: 'http://localhost:3000/products',
    lastModified: new Date(),
    changeFrequency: 'daily',
    priority: 0.9,
  },
  {
    url: 'http://localhost:3000/categories',
    lastModified: new Date(),
    changeFrequency: 'weekly',
    priority: 0.8,
  },
  {
    url: 'http://localhost:3000/deals',
    lastModified: new Date(),
    changeFrequency: 'daily',
    priority: 0.9,
  },
  {
    url: 'http://localhost:3000/search',
    lastModified: new Date(),
    changeFrequency: 'yearly',
    priority: 0.7,
  },
  {
    url: 'http://localhost:3000/login',
    lastModified: new Date(),
    changeFrequency: 'yearly',
    priority: 0.6,
  },
  {
    url: 'http://localhost:3000/register',
    lastModified: new Date(),
    changeFrequency: 'yearly',
    priority: 0.6,
  },
  {
    url: 'http://localhost:3000/account',
    lastModified: new Date(),
    changeFrequency: 'weekly',
    priority: 0.7,
  },
  {
    url: 'http://localhost:3000/account/orders',
    lastModified: new Date(),
    changeFrequency: 'weekly',
    priority: 0.6,
  },
  {
    url: 'http://localhost:3000/account/wishlist',
    lastModified: new Date(),
    changeFrequency: 'weekly',
    priority: 0.6,
  },
  {
    url: 'http://localhost:3000/account/addresses',
    lastModified: new Date(),
    changeFrequency: 'monthly',
    priority: 0.5,
  },
  {
    url: 'http://localhost:3000/account/payment-methods',
    lastModified: new Date(),
    changeFrequency: 'monthly',
    priority: 0.5,
  },
  {
    url: 'http://localhost:3000/cart',
    lastModified: new Date(),
    changeFrequency: 'daily',
    priority: 0.4,
  },
  {
    url: 'http://localhost:3000/checkout',
    lastModified: new Date(),
    changeFrequency: 'daily',
    priority: 0.3,
  },
  {
    url: 'http://localhost:3000/order-confirmation',
    lastModified: new Date(),
    changeFrequency: 'yearly',
    priority: 0.2,
  },
];

// Helper function to generate product URLs
async function generateProductUrls(): Promise<MetadataRoute.Sitemap> {
  try {
    const client = apolloClient;
    const { data } = await client.query({
      query: GET_PRODUCTS,
      variables: { limit: 1000 }, // Limit to avoid performance issues
    });

    return data.products.map((product: any) => ({
      url: `http://localhost:3000/products/${product.handle || product.id}`,
      lastModified: new Date(product.updatedAt),
      changeFrequency: 'weekly',
      priority: 0.8,
    }));
  } catch (error) {
    console.error('Error generating product URLs for sitemap:', error);
    return [];
  }
}

// Helper function to generate category URLs
async function generateCategoryUrls(): Promise<MetadataRoute.Sitemap> {
  try {
    const client = apolloClient;
    const { data } = await client.query({
      query: GET_CATEGORIES,
      variables: { limit: 100 }, // Limit to avoid performance issues
    });

    return data.categories.map((category: any) => ({
      url: `http://localhost:3000/categories/${category.handle || category.id}`,
      lastModified: new Date(category.updatedAt),
      changeFrequency: 'weekly',
      priority: 0.7,
    }));
  } catch (error) {
    console.error('Error generating category URLs for sitemap:', error);
    return [];
  }
}

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  // Generate dynamic URLs
  const productUrls = await generateProductUrls();
  const categoryUrls = await generateCategoryUrls();

  // Combine all URLs
  return [
    ...BASE_URLS,
    ...productUrls,
    ...categoryUrls,
  ];
}