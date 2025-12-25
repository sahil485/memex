// Wails Runtime JS bindings
export function Quit() {
  if (window.runtime && window.runtime.Quit) {
    window.runtime.Quit();
  }
}

export function WindowHide() {
  if (window.runtime && window.runtime.WindowHide) {
    window.runtime.WindowHide();
  }
}

export function EventsOn(eventName, callback) {
  if (window.runtime && window.runtime.EventsOn) {
    window.runtime.EventsOn(eventName, callback);
  }
}

export function EventsOff(eventName) {
  if (window.runtime && window.runtime.EventsOff) {
    window.runtime.EventsOff(eventName);
  }
}

export function EventsEmit(eventName, ...data) {
  if (window.runtime && window.runtime.EventsEmit) {
    window.runtime.EventsEmit(eventName, ...data);
  }
}

export function LogPrint(message) {
  if (window.runtime && window.runtime.LogPrint) {
    window.runtime.LogPrint(message);
  }
}

export function LogTrace(message) {
  if (window.runtime && window.runtime.LogTrace) {
    window.runtime.LogTrace(message);
  }
}

export function LogDebug(message) {
  if (window.runtime && window.runtime.LogDebug) {
    window.runtime.LogDebug(message);
  }
}

export function LogInfo(message) {
  if (window.runtime && window.runtime.LogInfo) {
    window.runtime.LogInfo(message);
  }
}

export function LogWarning(message) {
  if (window.runtime && window.runtime.LogWarning) {
    window.runtime.LogWarning(message);
  }
}

export function LogError(message) {
  if (window.runtime && window.runtime.LogError) {
    window.runtime.LogError(message);
  }
}

export function LogFatal(message) {
  if (window.runtime && window.runtime.LogFatal) {
    window.runtime.LogFatal(message);
  }
}
