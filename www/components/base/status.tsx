import React, { useEffect, useState } from 'react';

type Status = 'loading' | 'success' | 'error' ;

interface LoadingCircleProps {
  status?: Status;
}

const statusColors: Record<Status, { outer: string; inner: string }> = {
  loading: { outer: 'bg-blue-200', inner: 'bg-blue-500' },
  success: { outer: 'bg-green-200', inner: 'bg-green-500' },
  error: { outer: 'bg-red-200', inner: 'bg-red-500' },
};

const LoadingCircle: React.FC<LoadingCircleProps> = ({ status }) => {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    const intervalId = setInterval(() => {
      setIsVisible((prev) => !prev);
    }, 1000);

    return () => clearInterval(intervalId);
  }, []);

  let { outer, inner } = { outer: 'bg-gray-200', inner: 'bg-gray-500' }
  if (status) {
    const { outer: o, inner: i } = statusColors[status];
    outer = o;
    inner = i;
  }

  return (
    <div className="relative flex w-6 h-6">
      <div className={`absolute w-6 h-6 rounded-full ${outer} transition-opacity duration-500 ${isVisible ? 'opacity-100' : 'opacity-50'}`}>
        <div className={`absolute top-1 left-1 w-4 h-4 rounded-full ${inner} transition-opacity duration-500 ${isVisible ? 'opacity-100' : 'opacity-50'}`} />
      </div>
    </div>
  );
};

export default LoadingCircle;
