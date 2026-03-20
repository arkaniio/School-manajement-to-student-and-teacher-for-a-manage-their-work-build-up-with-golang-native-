"use client";

import React, { useEffect } from 'react';
import { Sidebar } from './Sidebar';
import { Navbar } from './Navbar';
import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from 'next/navigation';

export const DashboardLayout = ({ children }: { children: React.ReactNode }) => {
  const { user, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !user) {
      router.push('/login');
    }
  }, [user, isLoading, router]);

  if (isLoading) {
    return (
      <div className="h-screen w-full flex flex-col items-center justify-center"
        style={{ background: 'linear-gradient(135deg, #0a0e1a, #111827)' }}>
        <div className="relative mb-4">
          <div className="w-14 h-14 rounded-2xl flex items-center justify-center text-2xl mb-5 mx-auto"
            style={{ background: 'linear-gradient(135deg, #7c3aed, #4f46e5)', boxShadow: '0 8px 32px rgba(124,58,237,0.4)' }}>
            🏫
          </div>
          <div className="absolute -inset-1 rounded-2xl animate-pulse-ring"
            style={{ background: 'rgba(124,58,237,0.2)' }} />
        </div>
        <div className="flex gap-6">
          {[0, 1, 2].map(i => (
            <div key={i} className="w-2 h-2 rounded-full"
              style={{
                background: '#7c3aed',
                animation: `bounce 1.2s ease-in-out ${i * 0.2}s infinite`,
              }} />
          ))}
        </div>
        <style>{`@keyframes bounce { 0%,100%{transform:translateY(0)}50%{transform:translateY(-8px)} }`}</style>
      </div>
    );
  }

  if (!user) return null;

  return (
    <div className="h-screen w-full flex overflow-hidden" style={{ background: '#0a0e1a' }}>
      <Sidebar />
      <div className="flex-1 flex flex-col h-full min-w-0 overflow-hidden">
        <Navbar />
        <main className="flex-1 overflow-y-auto" style={{ background: '#0d1117' }}>
          {/* Ambient background orbs */}
          <div className="pointer-events-none fixed top-0 right-0 w-96 h-96 rounded-full opacity-30"
            style={{ background: 'radial-gradient(circle, rgba(124,58,237,0.08), transparent 70%)', transform: 'translate(30%, -30%)' }} />
          <div className="pointer-events-none fixed bottom-0 left-64 w-80 h-80 rounded-full opacity-20"
            style={{ background: 'radial-gradient(circle, rgba(59,130,246,0.08), transparent 70%)', transform: 'translate(-30%, 30%)' }} />
          
          <div className="relative max-w-7xl mx-auto px-6 py-8">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
};
