<template>
    <div>
        <manager :serviceProviderAccount="serviceProviderAccount" :dataConsumerAccount="dataConsumerAccount"></manager>
    </div>
</template>

<script>
    // @ is an alias to /src
    require("@/assets/css/style.css");
    export default {
        name: "admin",
        components: {
            manager: () => import("../components/manager.vue")
        },
        data: function () {
            return {
                serviceProviderAccount: undefined,
                dataConsumerAccount: undefined,
            }
        },
        computed: {
            account: function () {
                return this.$store.state.account
            }
        },
        created: function() {
            this.serviceProviderAccount = {
                address: "0x5eb76be14a2408508f64a2d35acf3b4d19198d45",
                privateKey: "275b9b0eb20cbf877c96ab88cbe11b7609229162c15f0f2d206d85492be94f1b"
            };
            this.dataConsumerAccount = {
                address: "0xd1A13B2C486d24D15BBf6b9376881EAB9AC9437D",
                privateKey: "A59D95FF309B945981D3D840BDC7C4E21520F0E26BB3D50A066ED6D845B2BA48"
            };
        }
    };
</script>
