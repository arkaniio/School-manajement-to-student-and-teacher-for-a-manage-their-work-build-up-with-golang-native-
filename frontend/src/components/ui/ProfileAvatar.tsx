"use client";

import React, { useState } from 'react';

interface ProfileAvatarProps {
  profileImage?: string | null;
  username?: string | null;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
}

const BASE_URL = 'http://localhost:8080/api/v1/image';

// Extract just the filename from a path like "uploadsUser\abc.jpg" or "uploadsUser/abc.jpg" → "abc.jpg"
function getImageUrl(profileImage: string): string {
  // Normalize Windows backslashes to forward slashes first
  const normalized = profileImage.replace(/\\/g, '/');
  const filename = normalized.split('/').pop() || profileImage;
  return `${BASE_URL}/${filename}`;
}

export const ProfileAvatar = ({ profileImage, username, className = '', size = 'md' }: ProfileAvatarProps) => {
  const [imgError, setImgError] = useState(false);

  const sizeClass = {
    sm: 'w-8 h-8 text-sm',
    md: 'w-9 h-9 text-sm',
    lg: 'w-16 h-16 text-2xl',
  }[size];

  const initial = username?.charAt(0).toUpperCase() || 'U';
  const imgUrl = profileImage && !imgError ? getImageUrl(profileImage) : null;

  return (
    <div className={`${sizeClass} rounded-full overflow-hidden flex items-center justify-center font-bold text-white shrink-0 ${className}`}
      style={!imgUrl ? { background: 'linear-gradient(135deg, #4f46e5, #7c3aed)' } : {}}>
      {imgUrl ? (
        <img
          src={imgUrl}
          alt={username || 'Profile'}
          className="w-full h-full object-cover"
          onError={() => setImgError(true)}
        />
      ) : (
        initial
      )}
    </div>
  );
};
