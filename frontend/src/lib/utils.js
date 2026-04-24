export function formatBytes(bytes) {
  if (bytes === 0) return "0 B";
  if (bytes < 0) return "0 B"; // Handle negative bytes as 0 B or however appropriate
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}
