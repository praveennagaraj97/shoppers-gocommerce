import axiosInstance from '@/api-routes';
import publicRoutes from '@/api-routes/publicRoutes';

export const loginAPI = async (email: string, password: string) => {
  const formData = new FormData();

  formData.append('email', email);
  formData.append('password', password);

  return axiosInstance.post(publicRoutes.Login, formData);
};
