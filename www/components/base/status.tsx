import React from 'react';

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
  let { outer, inner } = { outer: 'bg-gray-200', inner: 'bg-gray-500' }
  if (status) {
    const { outer: o, inner: i } = statusColors[status];
    outer = o;
    inner = i;
  }

  return (
    <div className="relative flex w-6 h-6">
      <div 
        className={`absolute w-6 h-6 rounded-full ${outer} animate-[ping_1.5s_ease-in-out_infinite]`}
        style={{ animationDelay: '0.2s' }}
      />
      <div 
        className={`absolute w-6 h-6 rounded-full ${outer} animate-[ping_1.5s_ease-in-out_infinite]`}
        style={{ animationDelay: '0.4s' }}
      />
      <div className={`absolute w-6 h-6 rounded-full ${outer}`}>
        <div 
          className={`absolute top-1 left-1 w-4 h-4 rounded-full ${inner} animate-[pulse_2s_cubic-bezier(0.4,0,0.6,1)_infinite]`} 
        />
      </div>
    </div>
  );
};

export default LoadingCircle;
