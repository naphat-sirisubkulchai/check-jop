"use client";

import { Toaster } from "@/components/ui/sonner";
import { ToasterProps } from "sonner";

/**
 * Custom Toaster component that extends the base Toaster from ui/sonner
 * - Override: Enhanced description text visibility (!text-gray-700 !opacity-100)
 * - Maintains: All base Toaster features and updates from ui/sonner
 *
 * This component wraps the base Toaster, so when ui/sonner updates,
 * this custom version will automatically inherit those updates.
 */
const CustomToaster = ({ ...props }: ToasterProps) => {
  return (
    <Toaster
      toastOptions={{
        classNames: {
          description: "!text-gray-700 !opacity-100",
        },
        ...props.toastOptions,
      }}
      {...props}
    />
  );
};

export { CustomToaster };
