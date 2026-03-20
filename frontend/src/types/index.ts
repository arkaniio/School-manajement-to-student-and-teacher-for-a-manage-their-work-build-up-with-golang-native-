export interface User {
  id: string;
  username: string;
  email: string;
  role: 'siswa' | 'guru' | 'admin';
  profile_image?: string;
}

export interface AuthResponse {
  token: string;
  refresh_token: string;
  user: User;
}

export interface Student {
  id: string;
  name: string;
  nisn: string;
  class: string;
  address?: string;
  user_id: string;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  deadline: string;
  file_url?: string;
  teacher_id: string;
}

export interface Absensi {
  id: string;
  student_id: string;
  status: 'h' | 'i' | 's' | 'a';
  date: string;
  notes?: string;
}
