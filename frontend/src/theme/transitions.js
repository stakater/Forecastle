// Transition and animation constants for Forecastle UI
// Provides consistent timing and easing across components

export const transitions = {
  duration: {
    shortest: 150,
    shorter: 200,
    short: 250,
    standard: 300,
    complex: 375,
    entering: 225,
    leaving: 195,
  },
  easing: {
    easeInOut: 'cubic-bezier(0.4, 0, 0.2, 1)',
    easeOut: 'cubic-bezier(0.0, 0, 0.2, 1)',
    easeIn: 'cubic-bezier(0.4, 0, 1, 1)',
    sharp: 'cubic-bezier(0.4, 0, 0.6, 1)',
  },
};

// Component-specific animation presets
export const animationPresets = {
  card: {
    hover: {
      transform: 'translateY(-2px)',
      transition: `all ${transitions.duration.shorter}ms ${transitions.easing.easeOut}`,
    },
    active: {
      transform: 'scale(0.98)',
      transition: `all ${transitions.duration.shortest}ms ${transitions.easing.easeIn}`,
    },
  },
  accordion: {
    expand: {
      transition: `all ${transitions.duration.standard}ms ${transitions.easing.easeInOut}`,
    },
  },
  fade: {
    enter: {
      transition: `opacity ${transitions.duration.entering}ms ${transitions.easing.easeOut}`,
    },
    exit: {
      transition: `opacity ${transitions.duration.leaving}ms ${transitions.easing.easeIn}`,
    },
  },
};
