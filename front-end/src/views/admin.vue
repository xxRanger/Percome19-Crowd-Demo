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
                address: "0xcb4dF655d08C9f054493809DCDb5a47F3A95eFAA",
                privateKey: "7a318fc20fbfe6a35681ca1ba5d8a74a07e54af99c9942a38bb1b03aad6173c1"
            };
        }
    };
</script>
