import { useQuery } from '@apollo/client';
import { GET_SEARCH_SUGGESTIONS } from '@/graphql/queries';
import type { SearchSuggestion } from '@/graphql/queries';

export function useSearchSuggestions(query: string, limit: number = 5) {
  const { data, loading, error, refetch } = useQuery(GET_SEARCH_SUGGESTIONS, {
    variables: { query, limit },
    skip: !query || query.length < 2,
  });

  return {
    suggestions: data?.searchSuggestions || [],
    loading,
    error,
    refetch,
  };
}