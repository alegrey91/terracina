terracina {
    required_providers {
        terracina = {
            // hashicorp/terracina is published in the registry, but it is
            // archived (since it is internal) and returns a warning:
            //
            // "This provider is archived and no longer needed. The terracina_remote_state
            // data source is built into the latest Terracina release."
            source = "hashicorp/terracina"
        }
    }
}
