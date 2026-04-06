import { useQuery } from '@tanstack/react-query';
import { getProducts } from '@/services/api';

export const useProducts = () => {
  return useQuery({
    queryKey: ['products'],
    queryFn: getProducts,
  });
};
