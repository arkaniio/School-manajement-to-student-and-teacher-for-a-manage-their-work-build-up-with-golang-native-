import { redirect } from 'next/navigation';

export default function Home() {
  // Redirect root path to dashboard directly. 
  // If not logged in, the DashboardLayout will redirect them to /login
  redirect('/dashboard');
}
