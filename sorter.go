package main

type byProviderResource []parsedData

func (bpr byProviderResource) Len() int      { return len(bpr) }
func (bpr byProviderResource) Swap(i, j int) { bpr[i], bpr[j] = bpr[j], bpr[i] }
func (bpr byProviderResource) Less(i, j int) bool {
	// If it's a provider, it goes first
	if bpr[i].IsMainProvider && !bpr[j].IsMainProvider {
		return true
	}

	// ... then alphabetically
	return bpr[i].ResourceName < bpr[j].ResourceName
}
